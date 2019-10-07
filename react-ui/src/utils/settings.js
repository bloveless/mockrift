export const getSettings = (env) => ({
    'API_URL': env === 'development' ? 'http://localhost:3499/admin/graphql' : '/admin/graphql'
});
