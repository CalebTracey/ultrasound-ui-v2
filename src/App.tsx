/* eslint-disable no-underscore-dangle */
import React, { FC, useEffect, useCallback, useState } from 'react'
import { useHistory, useLocation } from 'react-router-dom'
import { useAppSelector, useAppDispatch } from './redux/hooks'
import eventBus from './common/EventBus'
import './styles.scss'
import { logout, defineContentPath } from './redux/slices/auth'
import useItems from './hooks/useItems'
import { clearMessage } from './redux/slices/message'
import { IAppUser } from './schemas'
import useSubMenu from './hooks/useSubMenu'
import RoutesContainer from './containers/RoutesContainer'

const App: FC = () => {
    const { isAuth, loading, user, contentPath } = useAppSelector(
        (state) => state.auth
    )
    const [subId, setSubId] = useState('')
    const subMenuSelectedId = useAppSelector(
        (state) => state.subMenu.selected._id
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
    const isString = (value: unknown): value is string => {
        return !!value && !!(value as string)
    }
    const [itemResponse, getItems] = useItems({
        parentId: '',
        list: [],
        isLoading: false,
        error: null,
    })
    const [subMenuResponse, getSubMenu] = useSubMenu({
        id: subId,
        subMenuObj: {},
        isLoading: false,
        error: null,
    })

    useEffect(() => {
        if (isString(subMenuSelectedId)) {
            setSubId(subMenuSelectedId)
        }
    }, [subMenuSelectedId])

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

    // useEffect(() => {
    //     const ac = new AbortController()
    //     if (subId !== '') {
    //         eventBus.on('updateSubMenus', ac, () => {
    //             getSubMenu()
    //         })
    //     }
    //     return () => {
    //         ac.abort()
    //         eventBus.remove('updateSubMenus', getSubMenu)
    //     }
    // }, [getSubMenu, subMenuResponse, subId])

    // useEffect(() => {
    //     const ac = new AbortController()

    //     eventBus.on('updateItems', ac, () => {
    //         getItems()
    //     })
    //     return () => {
    //         ac.abort()
    //         eventBus.remove('updateItems', getItems)
    //     }
    // }, [getItems, itemResponse])

    return <RoutesContainer />
}
export default App
