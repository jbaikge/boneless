import { fetchUtils } from 'react-admin';
import { stringify } from 'query-string';

// Referencing environment variables from env.js in the public folder. This
// allows for compile-once and deploy-anywhere behavior
// https://www.freecodecamp.org/news/how-to-implement-runtime-environment-variables-with-create-react-app-docker-and-nginx-7f9d42a91d70/
const apiUrl = window._env_.API_URL;
const httpClient = fetchUtils.fetchJson;

// https://marmelab.com/react-admin/DataProviderWriting.html
const dataProvider = {
    /**
     * Search for resources
     *
     * @param {string} resource Resource name
     * @param {Object} params {
     *     pagination: {
     *         page: {int},
     *         perPage: {int}
     *     },
     *     sort: {
     *         field: {string},
     *         order: {string}
     *     },
     *     filter: {Object},
     *     meta: {Object}
     * }
     * @returns { data: {Record[]}, total: {int} }
     */
    getList: (resource, params) => {
        // https://otac0n.com/blog/2012/11/21/range-header-i-choose-you.html
        const { page, perPage } = params.pagination;
        const { field, order } = params.sort;
        const query = {
            sort:   `${field},${order}`,
            filter: JSON.stringify(params.filter),
        };
        const url = `${apiUrl}/${resource}?${stringify(query)}`;

        // TODO add options.headers with Range using page and perPage
        const range = [
            (page - 1) * perPage,
            page * perPage - 1,
        ];
        console.log(range);

        return httpClient(url).then(({ headers, json }) => ({
            data: json,
            total: json.length,
            // total: parseInt(headers.get('content-range').split('/').pop(), 10),
        }));
    },

    /**
     * Read a single resource, by id
     *
     * @param {string} resource Resource name
     * @param {Object} params { id: {mixed}, meta: {Object} }
     * @returns { data: {Record} }
     */
    getOne: (resource, params) => {
        const url = `${apiUrl}/${resource}/${params.id}`;

        return httpClient(url).then(({ json }) => ({ data: json }));
    },

    /**
     * Read a list of resource, by ids
     *
     * @param {string} resource Resource name
     * @param {Object} params { ids: {mixed[]}, meta: {Object} }
     * @returns { data: {Record[]} }
     */
    getMany: (resource, params) => {
        const query = {
            filter: JSON.stringify({ ids: params.ids }),
        };
        const url = `${apiUrl}/${resource}?${stringify(query)}`;
        return httpClient(url).then(({ json }) => ({ data: json }));
    },

    /**
     * Read a list of resources related to another one
     *
     * @param {string} resource Resource name
     * @param {Object} params {
     *     target: {string},
     *     id: {mixed},
     *     pagination: {
     *         page: {int},
     *         perPage: {int}
     *     },
     *     sort: {
     *         field: {string},
     *         order: {string}
     *     },
     *     filter: {Object},
     *     meta: {Object}
     * }
     * @returns { data: {Record[]}, total: {int} }
     */
    getManyReference: (resource, params) => {
        const { page, perPage } = params.pagination;
        const { field, order } = params.sort;
        const query = {
            sort: JSON.stringify([field, order]),
            range: JSON.stringify([(page - 1) * perPage, page * perPage - 1]),
            filter: JSON.stringify({
                ...params.filter,
                [params.target]: params.id,
            }),
        };
        const url = `${apiUrl}/${resource}?${stringify(query)}`;

        return httpClient(url).then(({ headers, json }) => ({
            data: json,
            total: parseInt(headers.get('content-range').split('/').pop(), 10),
        }));
    },

    /**
     * Create a single resource
     *
     * @param {string} resource Resource name
     * @param {Object} params { data: {Object}, meta: {Object} }
     * @returns { data: {Record} }
     */
    create: (resource, params) => {
        const url = `${apiUrl}/${resource}`;

        httpClient(url, {
            method: 'POST',
            body: JSON.stringify(params.data),
        }).then(({ json }) => ({
            data: { ...params.data, id: json.id },
        }))
    },

    /**
     * Updates a single resource
     *
     * @param {string} resource Resource name
     * @param {Object} params {
     *     id: {mixed},
     *     data: {Object},
     *     previousData: {Object},
     *     meta: {Object}
     * }
     * @returns { data: {Record} }
     */
    update: (resource, params) => {
        httpClient(`${apiUrl}/${resource}/${params.id}`, {
            method: 'PUT',
            body: JSON.stringify(params.data),
        }).then(({ json }) => ({ data: json }));
    },

    /**
     * Update multiple resources
     *
     * @param {string} resource Resource name
     * @param {Object} params {
     *     ids: {mixed[]},
     *     data: {Object},
     *     meta: {Object}
     * }
     * @returns { data: {mixed[]} } The ids which have been updated
     */
    updateMany: (resource, params) => {
        const query = {
            filter: JSON.stringify({ id: params.ids}),
        };
        const url = `${apiUrl}/${resource}?${stringify(query)}`;

        return httpClient(url, {
            method: 'PUT',
            body: JSON.stringify(params.data),
        }).then(({ json }) => ({ data: json }));
    },

    /**
     * Delete a single resource
     *
     * @param {string} resource Resource name
     * @param {Object} params {
     *     id: {mixed},
     *     previousData: {Object},
     *     meta: {Object}
     * }
     * @returns { data: {Record} } The record that has been deleted
     */
    delete: (resource, params) => {
        const url = `${apiUrl}/${resource}/${params.id}`

        httpClient(url, {
            method: 'DELETE',
        }).then(({ json }) => ({ data: json }));
    },

    /**
     * Delete multiple resources
     *
     * @param {string} resource Resource name
     * @param {Object} params { ids: {mixed[]}, meta: {Object} }
     * @returns { data: {mixed[]} } The ids of the deleted records (optional)
     */
    deleteMany: (resource, params) => {
        const query = {
            filter: JSON.stringify({ id: params.ids}),
        };
        const url = `${apiUrl}/${resource}?${stringify(query)}`;

        return httpClient(url, {
            method: 'DELETE',
            body: JSON.stringify(params.data),
        }).then(({ json }) => ({ data: json }));
    },
};

export default dataProvider;
