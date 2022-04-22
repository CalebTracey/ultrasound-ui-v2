import React, { useState, useEffect, FC } from 'react'
import { Media, Jumbotron, Container } from 'reactstrap'
import { AxiosResponse } from 'axios'
import LogoutButton from '../components/buttons/LogoutButton'
import LoginButton from '../components/login/LoginButton'
import RegisterButton from '../components/register/RegisterButton'
import DashboardButton from '../components/buttons/DashboardButton'
import { useAppSelector } from '../redux/hooks'
import { api } from '../service/api'

type TDate = {
    data: string
}
const Home: FC = () => {
    const { isAuth } = useAppSelector((state) => state.auth)
    const [content, setContent] = useState(null)

    const getDate = async () => {
        const date: Promise<AxiosResponse<TDate>> = api.get<TDate>(`date`)

        setContent((await date).data)
    }

    useEffect(() => {
        getDate()
    }, [])

    return (
        <>
            <header>
                <div className="button-wrapper">
                    {isAuth ? (
                        <>
                            <LogoutButton />
                            <DashboardButton />
                            <div className="date">{content}</div>
                        </>
                    ) : (
                        <>
                            <LoginButton />
                            <RegisterButton />
                            <div className="date">{content}</div>
                        </>
                    )}
                </div>
            </header>
            <div
                className="content"
                style={{
                    left: '0',
                    right: '0',
                    marginLeft: '1rem',
                    marginRight: '1rem',
                }}
            >
                <Jumbotron>
                    <Container>
                        <Media body>
                            <Media heading>
                                <div className="content___header">
                                    Emergency Ultrasound Training Application
                                </div>
                                <hr className="my-2" />
                            </Media>
                            <p className="content___text">
                                Please register or login to continue
                            </p>
                        </Media>
                    </Container>
                </Jumbotron>
            </div>
            <footer />
        </>
    )
}

export default Home
