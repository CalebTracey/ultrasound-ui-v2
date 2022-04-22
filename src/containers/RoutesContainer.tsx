import React, { useState, useEffect, useCallback, FC } from 'react'
import { UncontrolledAlert } from 'reactstrap'
import { useAppDispatch, useAppSelector } from '../redux/hooks'
import { clearMessage } from '../redux/slices/message'
import Routes from '../routes/Routes'

const RoutesContainer: FC = () => {
    const { text, error } = useAppSelector((state) => state.message)
    const [messageState, setMessageState] = useState(false)
    const [errorState, setErrorState] = useState(false)
    const [notification, setNotification] = useState(false)
    const dispatch = useAppDispatch()

    const messageToggle = useCallback(() => {
        setMessageState(!messageState)
    }, [messageState])
    const errorToggle = useCallback(() => {
        setErrorState(!errorState)
    }, [errorState])

    useEffect(() => {
        const timer = setTimeout(() => {
            if (messageState) {
                dispatch(clearMessage()).then(() => {
                    setNotification(false)
                    messageToggle()
                })
            }
            if (errorState) {
                dispatch(clearMessage()).then(() => {
                    setNotification(false)
                    errorToggle()
                })
            }
        }, 4000)
        return () => clearTimeout(timer)
    }, [messageToggle, errorToggle, dispatch, messageState, errorState])

    useEffect(() => {
        if (text && !error && !notification) {
            setNotification(true)
            messageToggle()
        }
        if (text && error && !notification) {
            setNotification(true)
            errorToggle()
        }
    }, [text, error, notification, errorToggle, messageToggle])
    return (
        <>
            <UncontrolledAlert
                style={{
                    display: 'relative',
                    zIndex: 100,
                    top: 0,
                    textAlign: 'center',
                }}
                isOpen={messageState}
                toggle={messageToggle}
                color="info"
            >
                <span className="span-text___Bold"> {text} </span>
            </UncontrolledAlert>

            <UncontrolledAlert
                style={{
                    display: 'relative',
                    zIndex: 100,
                    top: 0,
                    textAlign: 'center',
                }}
                isOpen={errorState}
                toggle={errorToggle}
                color="danger"
            >
                <span className="span-text___Bold"> {text} </span>
            </UncontrolledAlert>
            <Routes />
        </>
    )
}

export default RoutesContainer
