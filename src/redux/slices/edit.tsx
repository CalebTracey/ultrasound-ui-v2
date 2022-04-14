/* eslint-disable import/no-named-as-default */
/* eslint-disable @typescript-eslint/no-unused-vars */
/* eslint-disable no-param-reassign */
import { createAsyncThunk, createSlice } from '@reduxjs/toolkit'
import { AxiosResponse } from 'axios'
import { IListItem, IMessageResponse } from '../../schemas'
import { api } from '../../service/api'
import {
    getAllClassifications,
    resetClassificationSelection,
    editingClassification,
} from './classification'
import { removeListItem } from './subMenu'
import { removeItem } from './item'
import { newMessage } from './message'

type TDeleteItemPayload = {
    id: string
    type: 'subMenu' | 'classification'
    item: IListItem
}
type TDeleteDataPayload = {
    id: string
    type: 'subMenu' | 'classification'
}
type TDataNamePayload = { name: string }
type TDataNameProps = {
    id: string
    textValue: string
    type: 'subMenu' | 'classification'
}
type TItemNameProps = {
    id: string
    textValue: string
    item: IListItem
    type: 'subMenu' | 'classification'
}
interface editSliceState {
    loading: 'pending' | 'idle' | 'successful'
}

const initialSliceState: editSliceState = {
    loading: 'idle',
}

export const editDataName = createAsyncThunk(
    'edit/dataName',
    async (data: TDataNameProps, thunkApi) => {
        const { type, id, textValue } = data
        const newName = { name: textValue }
        const response = await api.post<TDataNamePayload, AxiosResponse>(
            `/edit/${type}/name/${id}`,
            newName
        )
        // thunkApi.dispatch(getAllClassifications())
        thunkApi.dispatch(newMessage(response.data.message))
    }
)

export const editItemName = createAsyncThunk(
    'edit/itemName',
    async (data: TItemNameProps, thunkApi) => {
        const { type, id, item, textValue } = data
        const newName = { newName: textValue, link: item.link, name: item.name }
        const response = await api.post<TDataNamePayload, AxiosResponse>(
            `/edit/${type}/item/name/${id}`,
            newName
        )
        // thunkApi.dispatch(getAllClassifications())
        thunkApi.dispatch(newMessage(response.data.message))
    }
)

export const deleteData = createAsyncThunk(
    'items/delete',
    async (data: TDeleteDataPayload, thunkApi) => {
        const { id, type } = data
        const response = await api.delete<TDataNamePayload, AxiosResponse>(
            `/delete-data/${type}/${id}`
        )
        thunkApi.dispatch(getAllClassifications())
        thunkApi.dispatch(resetClassificationSelection())
        thunkApi.dispatch(newMessage(response.data.message))
    }
)

export const deleteItem = createAsyncThunk(
    'items/delete',
    async (payload: TDeleteItemPayload, thunkApi) => {
        const { id, type, item } = payload
        api.post<IListItem, AxiosResponse>(
            `/delete-item/${type}/${id}`,
            item
        ).then((res: AxiosResponse<IMessageResponse>) => {
            thunkApi.dispatch(removeItem(item.name))
            thunkApi.dispatch(newMessage(res.data.message))
            if (type === 'subMenu') {
                thunkApi.dispatch(removeItem(item.link)) // remove from item slice
                thunkApi.dispatch(removeListItem(item.link)) // remove from subMenu slice
            }
            thunkApi.dispatch(editingClassification(true))
        })
    }
)

export const importData = createAsyncThunk(
    'edit/import',
    async (_, thunkApi) => {
        const response = await api.delete('/tables/clear').then(() => {
            thunkApi.dispatch(newMessage('Initializing the database'))
            api.post('/S3/update').then(
                (res: AxiosResponse<IMessageResponse>) => {
                    thunkApi.dispatch(getAllClassifications())
                    thunkApi.dispatch(newMessage(res.data.message))
                }
            )
        })
        return response
    }
)

export const updateData = createAsyncThunk('edit/update', async (_, thunkApi) =>
    api.post('/S3/synch').then((res: AxiosResponse<IMessageResponse>) => {
        thunkApi.dispatch(getAllClassifications())
        thunkApi.dispatch(newMessage(res.data.message))
    })
)

const editSlice = createSlice({
    name: 'edit',
    initialState: initialSliceState,
    reducers: {},
    extraReducers: (builder) => {
        builder.addCase(editDataName.pending, (state) => {
            state.loading = 'pending'
        })
        builder.addCase(editDataName.fulfilled, (state) => {
            state.loading = 'successful'
        })
        builder.addCase(deleteItem.pending, (state) => {
            state.loading = 'pending'
        })
        builder.addCase(deleteItem.fulfilled, (state) => {
            state.loading = 'successful'
        })
    },
})

export const editReducer = editSlice.reducer

export default editSlice.reducer
