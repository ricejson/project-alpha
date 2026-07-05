import { useCallback, useEffect, useState } from "react"
import { api } from "@/lib/api"
import type { Tag, TagPayload } from "@/types/tag"

interface UseTagsState {
  tags: Tag[]
  loading: boolean
  error: string | null
}

export function useTags() {
  const [state, setState] = useState<UseTagsState>({
    tags: [],
    loading: true,
    error: null,
  })

  const loadTags = useCallback(async () => {
    setState((current) => ({ ...current, loading: true, error: null }))
    try {
      const tags = await api.listTags()
      setState({ tags, loading: false, error: null })
    } catch (error) {
      setState((current) => ({
        ...current,
        loading: false,
        error: error instanceof Error ? error.message : "Failed to load tags",
      }))
    }
  }, [])

  useEffect(() => {
    void loadTags()
  }, [loadTags])

  const createTag = useCallback(
    async (payload: TagPayload) => {
      const tag = await api.createTag(payload)
      await loadTags()
      return tag
    },
    [loadTags],
  )

  const updateTag = useCallback(
    async (id: number, payload: TagPayload) => {
      const tag = await api.updateTag(id, payload)
      await loadTags()
      return tag
    },
    [loadTags],
  )

  const deleteTag = useCallback(
    async (id: number) => {
      await api.deleteTag(id)
      await loadTags()
    },
    [loadTags],
  )

  return {
    tags: state.tags,
    loading: state.loading,
    error: state.error,
    refresh: loadTags,
    createTag,
    updateTag,
    deleteTag,
  }
}
