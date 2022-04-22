/* eslint-disable no-underscore-dangle */
import React, { FC, useEffect, useCallback } from 'react'
import { useHistory, useLocation } from 'react-router-dom'
import { useAppSelector, useAppDispatch } from './redux/hooks'
import eventBus from './common/EventBus'
import './styles.scss'
import { logout, defineContentPath } from './redux/slices/auth'
import { clearMessage } from './redux/slices/message'
import { IAppUser } from './schemas'
import RoutesContainer from './containers/RoutesContainer'

const App: FC = () => {
    const { isAuth, loading, user, contentPath } = useAppSelector(
        (state) => state.auth
    )

    const dispatch = useAppDispatch()
    const history = useHistory()
    const location = useLocation()

    const isUser = (value: unknown): value is IAppUser => {
        return !!value && !!(value as IAppUser)
    }
    const isContentPath = (value: unknown): value is string => {
        return !!value && !!(value as string)
    }

    useEffect(() => {
        dispatch(clearMessage())
    }, [location.pathname, dispatch])

    const logOut = useCallback(() => {
        dispatch(logout())
    }, [dispatch])

    useEffect(() => {
        const ac = new AbortController()
        eventBus.on('logout', ac, () => {
            logOut()
            history.push('/home')
        })
        return () => {
            ac.abort()
            eventBus.remove('logout', logOut)
        }
    }, [dispatch, logOut, history])

    useEffect(() => {
        if (
            isAuth &&
            isUser(user) &&
            loading === 'successful' &&
            isContentPath(contentPath)
        ) {
            history.push(contentPath)
        }
        if (isUser(user)) {
            dispatch(defineContentPath(user))
        }
    }, [dispatch, history, isAuth, loading, contentPath, user])

    return <RoutesContainer />
}
export default App
