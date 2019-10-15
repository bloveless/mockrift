import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import { GET_APPS_QUERY } from '../../utils/queries';
import { useQuery } from '@apollo/react-hooks';

const Home = () => {
    const { loading, error, data } = useQuery(GET_APPS_QUERY);

    return (
        <div>
            Hello React,Webpack 4 & Babel 7!
            <ul>
                {loading && <li>Loading...</li>}
                {!loading && error && <li>There was an error: {error}</li>}
                {!loading && !error && data && data.apps && data.apps.map((app) => (
                    <li key={app.slug}>
                        <Link to={`/app/${app.slug}`}>{app.name || app.slug}</Link>
                    </li>
                ))}
            </ul>
        </div>
    );
};

export default Home;
