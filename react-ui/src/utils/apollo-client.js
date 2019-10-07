import ApolloClient from 'apollo-boost';

export const getClient = (uri) => {
    return new ApolloClient({ uri });
};
