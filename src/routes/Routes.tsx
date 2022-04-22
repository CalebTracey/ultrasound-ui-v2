import React, { FC, lazy, Suspense } from 'react'
import { Route, Switch } from 'react-router-dom'
import SyncLoader from 'react-spinners/SyncLoader'
import { useAppSelector } from '../redux/hooks'

const Dashboard = lazy(() => import('../containers/Dashboard'))
const Home = lazy(() => import('../containers/Home'))
const Login = lazy(() => import('../containers/Login'))
const ProtectedRoute = lazy(() => import('./ProtectedRoute'))
const Register = lazy(() => import('../containers/Register'))

// import Dashboard from '../containers/Dashboard'
// import Home from '../containers/Home'
// import Login from '../containers/Login'
// import ProtectedRoute from './ProtectedRoute'
// import Register from '../containers/Register'

const Routes: FC = () => {
    const { isAuth } = useAppSelector((state) => state.auth)

    return (
        <Suspense
            fallback={
                <div className="spinner">
                    <SyncLoader />
                </div>
            }
        >
            <Switch>
                <Route exact path={['/', '/home']} component={Home} />
                <Route exact path="/login" component={Login} />
                <Route exact path="/register" component={Register} />
                <ProtectedRoute
                    isAuthenticated={isAuth}
                    path="/dashboard"
                    authenticationPath="/login"
                    component={Dashboard}
                />
            </Switch>
        </Suspense>
    )
}

export default Routes
