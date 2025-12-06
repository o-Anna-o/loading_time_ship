// src/store/slices/filterSlice.ts
import { createSlice, PayloadAction } from '@reduxjs/toolkit'

interface FilterState {
  search: string
  appliedSearch: string
}

const initialState: FilterState = {
  search: '',
  appliedSearch: ''
}

const filterSlice = createSlice({
  name: 'filter',
  initialState,
  reducers: {
    setSearch(state, action: PayloadAction<string>) {
      state.search = action.payload
    },
    applySearch(state) {
      state.appliedSearch = state.search
    }
  }
})

export const { setSearch, applySearch } = filterSlice.actions
export default filterSlice.reducer
