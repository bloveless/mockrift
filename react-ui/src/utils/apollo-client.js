import ApolloClient from 'apollo-boost';
import { getSettings } from "./settings";

export const getDefaultClient = () => {
    return getClient(getSettings(process.env.MODE).API_URL);
};

export const getClient = (uri) => {
    return new ApolloClient({ uri });
};
