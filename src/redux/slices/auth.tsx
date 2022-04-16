/* eslint-disable no-param-reassign */
import {
    createAsyncThunk,
    createSlice,
    PayloadAction,
    Reducer,
} from '@reduxjs/toolkit'
import { AxiosResponse } from 'axios'
import AuthService from '../../service/auth-service'
import TokenService from '../../service/token-service'
import { IAppUser } from '../../schemas'
import { api } from '../../service/api'

type TLogin = { username: string; password: string }
type TRegister = {
    fullName: string
    username: string
    email: string
    password: string
    confirmPassword: string
}
interface TUserLoginResponse extends AxiosResponse {
    accessToken: string
    email: string
    id: string
    roles: string[]
    tokenType: string
}
interface authSliceState {
    isAuth: boolean
    user: IAppUser | Record<string, null>
    loading: 'idle' | 'pending' | 'successful'
    error: string | null
    contentPath: '/dashboard' | '/dashboard/admin' | null
    showEdit: boolean
}

const data: string | null = localStorage.getItem('user')
const user = data ? JSON.parse(data) : null

const initialAuthState: authSliceState = user
    ? {
          isAuth: true,
          user,
          loading: 'successful',
          error: null,
          contentPath: null,
          showEdit: false,
      }
    : {
          isAuth: false,
          user: {},
          loading: 'idle',
          error: null,
          contentPath: null,
          showEdit: false,
      }

const isUser = (value: unknown): value is IAppUser => {
    return !!value && !!(value as IAppUser)
}

export const login = createAsyncThunk(
    'auth/login',
    async (credentials: TLogin) => {
        const response = await api.post(`auth/sign-in`, credentials)
        return response.data
    }
)

export const logout = createAsyncThunk('auth/logout', async () => {
    return AuthService.logoutService()
})

// export const userRegister = createAsyncThunk(
//     'auth/register',
//     async (regCredentials: TRegister) => {
//         const response = await api.post<TRegister, TUserLoginResponse>(
//             `auth/sign-up`,
//             regCredentials
//         )
//         TokenService.setUser(response.data)
//         return response.data
//     }
// )

export const userRegister = createAsyncThunk(
    'auth/register',
    async (regCredentials: TRegister) => {
        const response = await api.post<TRegister>(
            `auth/sign-up`,
            regCredentials
        )

        return response.data
    }
)

export const authSlice = createSlice({
    name: 'auth',
    initialState: initialAuthState,
    reducers: {
        defineContentPath: (state, action: PayloadAction<IAppUser>) => {
            const userData = action.payload
            if (isUser(userData) && userData.roles !== undefined) {
                if (userData.roles.includes('ROLE_ADMIN')) {
                    state.contentPath = '/dashboard/admin'
                } else {
                    state.contentPath = '/dashboard'
                }
            } else {
                state.isAuth = false
            }
        },
        registerSuccess: (state, action: PayloadAction<IAppUser>) => {
            const userDetails = action.payload
            state.user = userDetails
            state.isAuth = true
            state.isAuth = false
        },
        registerFail: (state) => {
            state.isAuth = false
        },
        userRefreshToken: (state, action: PayloadAction<string>) => {
            const token = action.payload
            state.user = { ...user, accessToken: token }
        },
        showEditToggle: (state) => {
            const currentState = state.showEdit
            state.showEdit = !currentState
        },
    },
    extraReducers: (builder) => {
        builder.addCase(login.pending, (state) => {
            state.loading = 'pending'
        })
        builder.addCase(login.fulfilled, (state, action) => {
            const userData = action.payload
            if (isUser(userData)) {
                TokenService.setUser(userData)
                state.user = userData
                state.isAuth = true
                state.loading = 'successful'
                state.error = ''
                if (userData.roles.includes('ROLE_ADMIN')) {
                    state.contentPath = '/dashboard/admin'
                }
            } else {
                state.isAuth = false
                state.loading = 'idle'
                state.user = {}
                state.error = 'Login failed'
            }
        })
        builder.addCase(login.rejected, (state) => {
            state.isAuth = false
            state.user = {}
            state.loading = 'idle'
            state.error = 'Login failed - try again'
        })
        builder.addCase(userRegister.fulfilled, (state, action) => {
            state.user = {}
            state.isAuth = false
            state.loading = 'successful'
            state.error = 'Registration received. Account pending.'

            // const userData = action.payload
            // if (isUser(userData)) {
            //     state.user = userData
            //     state.isAuth = true
            //     state.loading = 'successful'
            //     state.error = ''
            //     if (userData.roles.includes('ROLE_ADMIN')) {
            //         state.contentPath = '/dashboard/admin'
            //     }
            // } else {
            //     state.isAuth = false
            //     state.loading = 'idle'
            //     state.user = {}
            //     state.error = 'Registration failed'
            // }
        })
        builder.addCase(userRegister.rejected, (state) => {
            state.isAuth = false
            state.loading = 'idle'
            state.user = {}
            state.error = 'Registration failed - username or email in use'
        })
        builder.addCase(logout.fulfilled, (state) => {
            state.isAuth = false
            state.loading = 'idle'
            state.error = ''
            state.user = {}
        })
    },
})
export const {
    showEditToggle,
    registerSuccess,
    registerFail,
    defineContentPath,
    userRefreshToken,
} = authSlice.actions

export default authSlice.reducer as Reducer<typeof initialAuthState>
