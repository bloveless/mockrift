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

const UPDATE_REQUEST_BODY = gql`
    mutation modifyRequestBody($id: ID!, $body: String!) {
        modifyRequestBody(id: $id, body: $body) {
            id
            url
            method
            body
        }
    }
`;

export {
    GET_APPS_QUERY,
    GET_APP_QUERY,
    UPDATE_REQUEST_BODY,
};
