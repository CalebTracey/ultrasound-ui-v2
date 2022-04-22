type TUserLoginResponse = {
    accessToken: string
    email: string
    id: string
    roles: string[]
    tokenType: string
}

const getLocalRefreshToken = (): string | undefined => {
    const userStorage = localStorage.getItem('user')
    if (userStorage) {
        const user = JSON.parse(userStorage)
        return user.refreshToken
    }
    return undefined
}

const getLocalAccessToken = (): string | undefined => {
    const userStorage = localStorage.getItem('user')
    if (userStorage) {
        const user = JSON.parse(userStorage)
        return user.accessToken
    }
    return undefined
}

const setUser = (user: TUserLoginResponse): void => {
    localStorage.setItem('user', JSON.stringify(user))
}

const removeUser = (): void => {
    localStorage.removeItem('user')
}

const TokenService = {
    getLocalRefreshToken,
    getLocalAccessToken,
    setUser,
    removeUser,
}

export default TokenService
