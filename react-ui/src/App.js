import React from 'react';
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom';
import { ApolloProvider } from '@apollo/react-hooks';
import Home from './pages/Home';
import AppDetails from './pages/AppDetails';
import { getClient } from './utils/apollo-client';
import { getSettings } from './utils/settings';

const App = () => (
    <ApolloProvider client={getClient(getSettings(process.env.MODE).API_URL)}>
        <Router basename="/admin">
            <Switch>
                <Route path="/app/:slug" component={AppDetails}/>
                <Route path="/" component={Home}/>
            </Switch>
        </Router>
    </ApolloProvider>
);

export default App;
