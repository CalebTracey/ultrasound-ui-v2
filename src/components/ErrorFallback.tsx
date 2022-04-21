import React, { FC } from 'react'

type Props = {
    err: Error
}

const ErrorFallback: FC<Props> = ({ err }): JSX.Element => {
    return (
        <div role="alert">
            <p>Something went wrong:</p>
            <pre>{err.message}</pre>
        </div>
    )
}

export default ErrorFallback
