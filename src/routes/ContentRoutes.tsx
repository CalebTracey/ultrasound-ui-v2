/* eslint-disable react/prop-types */
import React, { FC } from 'react'
import { Route, Switch, Redirect } from 'react-router-dom'

import EditObject from '../containers/EditObject'
import VideoPlayer from '../components/content/VideoPlayer'
import Classification from '../containers/Classification'
import ContentHome from '../components/content/ContentHome'

interface Props {
    routePath: string
}

const ContentRoutes: FC<Props> = ({ routePath }) => {
    return (
        <Switch>
            <Route path={`${routePath}/home`} exact component={ContentHome} />
            <Route
                path={`${routePath}/classification/:id`}
                component={Classification}
            />
            <Route path="/dashboard/admin/move/:id" component={EditObject} />
            <Route path={`${routePath}/video/:id`} component={VideoPlayer} />
            <Redirect from={`${routePath}`} to={`${routePath}/home`} exact />
        </Switch>
    )
}

export default ContentRoutes
