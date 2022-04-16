/* eslint-disable jsx-a11y/label-has-associated-control */
/* eslint-disable react/prop-types */
import React, { FC, useEffect, useState, useCallback } from 'react'
import Switch from 'react-switch'
import { AxiosError, AxiosResponse } from 'axios'
import Logout from '../buttons/LogoutButton'
import UserInfoHeader from '../UserInfoHeader'
import { useAppSelector, useAppDispatch } from '../../redux/hooks'
import { importData, updateData } from '../../redux/slices/edit'
import { newError } from '../../redux/slices/message'
import { api } from '../../service/api'
import { IAppUser } from '../../schemas'
import { userRegister } from '../../redux/slices/auth'
import WarningModal from '../WarningModal'
import DatabaseDropdown from '../DatabaseDropdown'

const Header: FC = () => {
    const { auth } = useAppSelector((state) => state)
    const { user } = auth
    const [content, setContent] = useState<string | null>(null)
    const [initModal, setInitModal] = useState(false)
    const [updateModal, setUpdateModal] = useState(false)
    const [dropDownOpen, setDropdownOpen] = useState(false)
    // const [switchChecked, setSwitchChecked] = useState(false)
    const isAdmin = user.roles?.includes('ROLE_ADMIN')
    const dispatch = useAppDispatch()
    const isUser = (value: unknown): value is IAppUser => {
        return !!value && !!(value as IAppUser)
    }
    // const handleSwitchChange = () => {
    //     dispatch(showEditToggle())
    //     setSwitchChecked(!switchChecked)
    // }
    const dropDownToggle = () => setDropdownOpen(!dropDownOpen)

    const databaseInitToggle = useCallback(() => {
        setInitModal(!initModal)
    }, [initModal])

    const databaseUpdateToggle = useCallback(() => {
        setUpdateModal(!updateModal)
    }, [updateModal])

    const handleDatabaseInit = () => {
        dispatch(importData()).catch((res: AxiosResponse<AxiosError>) => {
            dispatch(newError(res.data.message))
        })
        databaseInitToggle()
    }
    const handleDatabaseUpdate = () => {
        dispatch(updateData())
        databaseUpdateToggle()
    }
    useEffect(() => {
        const getDate = async () => {
            const date = await api.get(`date`)
            setContent(date.data)
        }
        getDate()
    }, [])

    // useEffect(() => {
    //     if (switchChecked !== showEdit) {
    //         setSwitchChecked(showEdit)
    //     }
    // }, [switchChecked, showEdit])

    return (
        <>
            <header style={{ zIndex: 10 }}>
                <div className="button-wrapper">
                    <div className="date">{content}</div>
                    {isAdmin && (
                        <div
                            className="form-group"
                            style={{ marginLeft: '1rem' }}
                        >
                            <DatabaseDropdown
                                dropDownOpen={dropDownOpen}
                                dropDownToggle={dropDownToggle}
                                // databaseInitToggle={databaseInitToggle}
                                databaseUpdateToggle={databaseUpdateToggle}
                            />
                        </div>
                    )}

                    <Logout />
                    {user.email && isUser(userRegister) && (
                        <UserInfoHeader email={user.email} />
                    )}
                    {/* {isAdmin && (
                        <div
                            style={{
                                position: 'absolute',
                                left: '1rem',
                            }}
                        >
                            <label>
                                <Switch
                                    height={14}
                                    width={28}
                                    handleDiameter={12}
                                    onChange={handleSwitchChange}
                                    checked={switchChecked}
                                />
                            </label>
                        </div>
                    )} */}
                    {/* <SearchBar /> */}
                </div>
            </header>
            <WarningModal
                actionText="Reset the "
                itemText="database"
                setModal={initModal}
                toggleCallback={databaseInitToggle}
                modalAction={handleDatabaseInit}
            />
            <WarningModal
                actionText="Update the "
                itemText="database"
                setModal={updateModal}
                toggleCallback={databaseUpdateToggle}
                modalAction={handleDatabaseUpdate}
            />
        </>
    )
}

export default Header
