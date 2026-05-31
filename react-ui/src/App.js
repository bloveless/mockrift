import React from 'react';
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom';
import { ApolloProvider } from '@apollo/react-hooks';
import Home from './pages/Home';
import AppDetails from './pages/AppDetails';
import { getDefaultClient } from './utils/apollo-client';

const App = () => (
    <ApolloProvider client={getDefaultClient()}>
        <Router basename="/admin">
            <Switch>
                <Route path="/app/:slug" component={AppDetails}/>
                <Route path="/" component={Home}/>
            </Switch>
        </Router>
    </ApolloProvider>
);

export default App;
