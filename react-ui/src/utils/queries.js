import { gql } from 'apollo-boost';

const GET_APPS_QUERY = gql`
    query getApps {
        apps {
            slug
            name
        }
    }
`;

const GET_APP_QUERY = gql`
    query getApp($slug: String!) {
        app(slug: $slug) {
            name
            slug
            requests {
                id
                url
                method
                header {
                    name
                    value
                }
                body
                responses {
                    id 
                    active
                    body
                    status_code
                    header {
                        name
                        value
                    }
                }
            }
        }
    }
`;

export { GET_APPS_QUERY, GET_APP_QUERY };
