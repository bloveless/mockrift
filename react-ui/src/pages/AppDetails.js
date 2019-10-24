import React, { useState } from 'react';
import { useQuery } from '@apollo/react-hooks';
import { GET_APP_QUERY } from '../utils/queries';
import RequestDetails from '../components/RequestDetails';

const AppDetails = ({ match }) => {
    const { loading, error, data } = useQuery(GET_APP_QUERY, {
        variables: {
            slug: match.params.slug
        }
    });
    const app = data ? data.app : null;

    if (!loading && error) {
        return <div>An error occurred {error}</div>
    }

    if (loading) {
        return <div>Loading...</div>
    }

    return <div>
        <h1>App: {app.name}</h1>
        <div>{app.slug}</div>
        <RequestDetails app={app}/>
    </div>;
};

export default AppDetails;
