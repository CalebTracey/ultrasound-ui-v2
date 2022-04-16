import { api } from './api'
import TokenService from './token-service'

export type TUserDetails = {
    fullName: string
    username: string
    email: string
}
type TNewUser = {
    fullname: string
    username: string
    email: string
    password: string
}
type TNewUserResponse = {
    message: string
}
type TUserLogin = {
    username: string
    password: string
}
type TUserLoginResponse = {
    accessToken: string
    email: string
    id: string
    roles: string[]
    tokenType: string
}
const { removeUser, setUser } = TokenService

const userDetails = async (username: string): Promise<TUserDetails> => {
    const response = await api.get(`auth/user/${username}`)
    return response.data
}

const registerService = async (data: TNewUser): Promise<TNewUserResponse> =>
    api.post<TNewUser, TNewUserResponse>(`auth/sign-up`, data)

const loginService = async (
    details: TUserLogin
): Promise<TUserLoginResponse> => {
    const response = await api
        .post<TUserLogin, TUserLoginResponse>(`auth/sign-in`, details)
        .then((res) => {
            if (res) {
                setUser(res)
            }
            return res
        })
    return response
}

const logoutService = (): void => {
    removeUser()
}

const AuthService = {
    userDetails,
    registerService,
    loginService,
    logoutService,
}

export default AuthService
