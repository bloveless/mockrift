import React, { useState } from 'react';
import { useQuery } from '@apollo/react-hooks';
import { GET_APP_QUERY } from '../../utils/queries';
import RequestDetails from './components/RequestDetails';

const AppDetails = ({ match }) => {
    const { loading, error, data: { app } } = useQuery(GET_APP_QUERY, {
        variables: {
            slug: match.params.slug
        }
    });

    if (!loading && error) {
        return <div>An error occurred</div>
    }

    if (loading) {
        return <div>Loading...</div>
    }

    return <div>
        <div>{app.name}</div>
        <div>{app.slug}</div>
        <RequestDetails app={app}/>
    </div>;
};

export default AppDetails;
