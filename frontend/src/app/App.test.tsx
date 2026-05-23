import { fireEvent, render, screen, waitFor, within } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { beforeEach, describe, expect, it, vi } from "vitest"
import { App } from "@/app/App"
import type { Tag } from "@/types/tag"
import type { Ticket } from "@/types/ticket"

const now = "2026-05-23T04:00:00Z"

let tags: Tag[]
let tickets: Ticket[]
let nextTagId: number
let nextTicketId: number

function resetData() {
  tags = [
    { id: 1, name: "backend", color: "#2563eb", createdAt: now, updatedAt: now },
    { id: 2, name: "frontend", color: "#16a34a", createdAt: now, updatedAt: now },
  ]
  tickets = [
    {
      id: 1,
      title: "Build API",
      description: "GORM and Gin endpoints",
      completed: false,
      completedAt: null,
      tags: [tags[0]],
      createdAt: now,
      updatedAt: now,
    },
    {
      id: 2,
      title: "Ship UI",
      description: "React ticket workflow",
      completed: true,
      completedAt: now,
      tags: [tags[1]],
      createdAt: now,
      updatedAt: now,
    },
  ]
  nextTagId = 3
  nextTicketId = 3
}

beforeEach(() => {
  resetData()
  vi.restoreAllMocks()
  vi.spyOn(globalThis, "fetch").mockImplementation(mockFetch)
})

describe("App", () => {
  it("渲染 ticket 列表", async () => {
    render(<App />)

    expect(await screen.findByText("Build API")).toBeInTheDocument()
    expect(screen.getByText("Ship UI")).toBeInTheDocument()
    expect(screen.getByText("backend")).toBeInTheDocument()
  })

  it("按标题搜索 ticket", async () => {
    const user = userEvent.setup()
    render(<App />)

    await screen.findByText("Build API")
    await user.type(screen.getByLabelText("搜索标题"), "Ship")

    await waitFor(() => {
      expect(screen.queryByText("Build API")).not.toBeInTheDocument()
      expect(screen.getByText("Ship UI")).toBeInTheDocument()
    })
  })

  it("按标签筛选 ticket", async () => {
    const user = userEvent.setup()
    render(<App />)

    await screen.findByText("Build API")
    await user.click(screen.getByRole("button", { name: /标签筛选/ }))
    await user.click(screen.getByRole("button", { name: /frontend/ }))

    await waitFor(() => {
      expect(screen.queryByText("Build API")).not.toBeInTheDocument()
      expect(screen.getByText("Ship UI")).toBeInTheDocument()
    })
  })

  it("校验并创建 ticket", async () => {
    const user = userEvent.setup()
    render(<App />)

    await screen.findByText("Build API")
    await user.click(screen.getByRole("button", { name: "新建 Ticket" }))

    const dialog = await screen.findByRole("dialog", { name: "新建 Ticket" })
    await user.click(within(dialog).getByRole("button", { name: "保存" }))
    expect(await within(dialog).findByText("标题不能为空")).toBeInTheDocument()

    await user.type(within(dialog).getByLabelText("标题"), "Write tests")
    await user.click(within(dialog).getByRole("button", { name: "保存" }))

    expect(await screen.findByText("Write tests")).toBeInTheDocument()
  })

  it("删除 ticket 前需要确认", async () => {
    const user = userEvent.setup()
    render(<App />)

    await screen.findByText("Build API")
    fireEvent.click(screen.getByRole("button", { name: "删除 Build API" }))

    const dialog = await screen.findByRole("dialog", { name: "删除 Ticket" })
    expect(within(dialog).getByText("确认删除「Build API」？")).toBeInTheDocument()
    await user.click(within(dialog).getByRole("button", { name: "删除" }))

    await waitFor(() => {
      expect(screen.queryByText("Build API")).not.toBeInTheDocument()
    })
  })

  it("展示 API 错误", async () => {
    const user = userEvent.setup()
    render(<App />)

    await screen.findByText("Build API")
    await user.click(screen.getByRole("button", { name: "标签" }))

    const manager = screen.getByRole("dialog")
    await user.type(within(manager).getByLabelText("标签名称"), "backend")
    await user.click(within(manager).getByRole("button", { name: "创建" }))

    expect(await within(manager).findByText("tag name already exists")).toBeInTheDocument()
  })
})

async function mockFetch(input: RequestInfo | URL, init?: RequestInit): Promise<Response> {
  const url = new URL(typeof input === "string" || input instanceof URL ? input.toString() : input.url)
  const path = url.pathname.replace("/api/v1", "")
  const method = init?.method ?? "GET"

  if (path === "/tags") {
    if (method === "GET") return jsonResponse([...tags])
    if (method === "POST") {
      const payload = parseBody<{ name: string; color: string | null }>(init)
      if (tags.some((tag) => tag.name.toLowerCase() === payload.name.toLowerCase())) {
        return errorResponse(409, "CONFLICT", "tag name already exists")
      }
      const tag = { id: nextTagId++, name: payload.name, color: payload.color, createdAt: now, updatedAt: now }
      tags.push(tag)
      return jsonResponse(tag, 201)
    }
  }

  const tagMatch = path.match(/^\/tags\/(\d+)$/)
  if (tagMatch) {
    const id = Number(tagMatch[1])
    if (method === "PUT") {
      const payload = parseBody<{ name: string; color: string | null }>(init)
      const tag = tags.find((item) => item.id === id)
      if (!tag) return errorResponse(404, "NOT_FOUND", "tag not found")
      Object.assign(tag, payload, { updatedAt: now })
      return jsonResponse(tag)
    }
    if (method === "DELETE") {
      tags = tags.filter((tag) => tag.id !== id)
      tickets = tickets.map((ticket) => ({ ...ticket, tags: ticket.tags.filter((tag) => tag.id !== id) }))
      return new Response(null, { status: 204 })
    }
  }

  if (path === "/tickets") {
    if (method === "GET") return jsonResponse(listTickets(url))
    if (method === "POST") {
      const payload = parseBody<{ title: string; description: string; tagIds: number[] }>(init)
      const ticketTags = tags.filter((tag) => payload.tagIds.includes(tag.id))
      const ticket: Ticket = {
        id: nextTicketId++,
        title: payload.title,
        description: payload.description,
        completed: false,
        completedAt: null,
        tags: ticketTags,
        createdAt: now,
        updatedAt: now,
      }
      tickets = [ticket, ...tickets]
      return jsonResponse(ticket, 201)
    }
  }

  const ticketMatch = path.match(/^\/tickets\/(\d+)$/)
  if (ticketMatch) {
    const id = Number(ticketMatch[1])
    if (method === "PUT") {
      const payload = parseBody<{ title: string; description: string; tagIds: number[] }>(init)
      tickets = tickets.map((ticket) =>
        ticket.id === id
          ? {
              ...ticket,
              title: payload.title,
              description: payload.description,
              tags: tags.filter((tag) => payload.tagIds.includes(tag.id)),
              updatedAt: now,
            }
          : ticket,
      )
      return jsonResponse(tickets.find((ticket) => ticket.id === id))
    }
    if (method === "DELETE") {
      tickets = tickets.filter((ticket) => ticket.id !== id)
      return new Response(null, { status: 204 })
    }
  }

  const completeMatch = path.match(/^\/tickets\/(\d+)\/(complete|uncomplete)$/)
  if (completeMatch && method === "POST") {
    const id = Number(completeMatch[1])
    const completed = completeMatch[2] === "complete"
    tickets = tickets.map((ticket) =>
      ticket.id === id ? { ...ticket, completed, completedAt: completed ? now : null, updatedAt: now } : ticket,
    )
    return jsonResponse(tickets.find((ticket) => ticket.id === id))
  }

  return errorResponse(404, "NOT_FOUND", `unhandled ${method} ${path}`)
}

function listTickets(url: URL) {
  const title = url.searchParams.get("title")?.toLowerCase() ?? ""
  const completed = url.searchParams.get("completed")
  const tagIds = (url.searchParams.get("tagIds") ?? "")
    .split(",")
    .filter(Boolean)
    .map((value) => Number(value))
  const page = Number(url.searchParams.get("page") ?? "1")
  const pageSize = Number(url.searchParams.get("pageSize") ?? "20")

  let items = [...tickets]
  if (title) items = items.filter((ticket) => ticket.title.toLowerCase().includes(title))
  if (completed === "true") items = items.filter((ticket) => ticket.completed)
  if (completed === "false") items = items.filter((ticket) => !ticket.completed)
  if (tagIds.length) {
    items = items.filter((ticket) => tagIds.every((tagId) => ticket.tags.some((tag) => tag.id === tagId)))
  }

  const total = items.length
  const start = (page - 1) * pageSize
  return {
    items: items.slice(start, start + pageSize),
    pagination: {
      page,
      pageSize,
      total,
      totalPages: total ? Math.ceil(total / pageSize) : 0,
    },
  }
}

function parseBody<T>(init?: RequestInit): T {
  return JSON.parse(String(init?.body ?? "{}")) as T
}

function jsonResponse<T>(data: T, status = 200) {
  return new Response(JSON.stringify({ data, error: null }), {
    status,
    headers: { "Content-Type": "application/json" },
  })
}

function errorResponse(status: number, code: string, message: string) {
  return new Response(JSON.stringify({ data: null, error: { code, message } }), {
    status,
    headers: { "Content-Type": "application/json" },
  })
}
