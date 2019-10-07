import { gql } from 'apollo-boost';

const GET_APPS_QUERY = gql`
    query getApps {
        apps {
            slug
            name
        }
    }
`;

export { GET_APPS_QUERY };
