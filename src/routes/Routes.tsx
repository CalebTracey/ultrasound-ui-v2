import React, { FC, lazy, Suspense } from 'react'
import { Route } from 'react-router-dom'
import SyncLoader from 'react-spinners/SyncLoader'
// import { IAppUser } from '../schemas'
// import Dashboard from '../containers/Dashboard'
import { useAppSelector } from '../redux/hooks'

const Dashboard = lazy(() => import('../containers/Dashboard'))
const Home = lazy(() => import('../containers/Home'))
const Login = lazy(() => import('../containers/Login'))
const ProtectedRoute = lazy(() => import('./ProtectedRoute'))
const Register = lazy(() => import('../containers/Register'))
const Routes: FC = () => {
    const { isAuth } = useAppSelector((state) => state.auth)

    // const isUser = (value: unknown): value is IAppUser => {
    //     return !!value && !!(value as IAppUser)
    // }
    // const isAdmin = isUser(user) && user.roles?.includes('ROLE_ADMIN')

    return (
        <Suspense
            fallback={
                <div className="spinner">
                    <SyncLoader />
                </div>
            }
        >
            <Route exact path={['/', '/home']} component={Home} />
            <Route exact path="/login" component={Login} />
            <Route exact path="/register" component={Register} />
            {/* <Switch> */}
            <ProtectedRoute
                isAuthenticated={isAuth}
                path="/dashboard"
                authenticationPath="/login"
                component={Dashboard}
            />
            {/* <ProtectedRoute
                    isAuthenticated={isAdmin}
                    path="/admin/dashboard"
                    authenticationPath="/login"
                    component={Dashboard}
                />
                )
            </Switch> */}
        </Suspense>
    )
}

export default Routes
