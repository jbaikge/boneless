import simpleRestProvider from 'ra-data-simple-rest';
import { fetchUtils } from 'ra-core';

const API_URL = process.env.REACT_APP_API_URL;
const baseDataProvider = simpleRestProvider(API_URL);

const uploadFile = (key, file) => {
  let location = '';
  return fetchUtils.fetchJson(`${API_URL}/files/url`, {
    method: 'POST',
    body: JSON.stringify({
      key: key,
      expires: '5m',
      content_type: file.type,
    }),
  }).then(({ json }) => new Promise((resolve, reject) => {
    location = json.location;
    const reader = new FileReader();
    reader.addEventListener('load', () => resolve({
      url: json.url,
      method: json.method,
      body: reader.result
    }));
    reader.addEventListener('error', reject);
    reader.readAsBinaryString(file);
  })).then(({ url, method, body }) => fetchUtils.fetchJson(url, {
    method: method,
    headers: new Headers({ 'Content-Type': file.type }),
    body: body,
  })).then(() => location);
}

const dataProvider = {
  ...baseDataProvider,
  update: (resource, params) => {
    // No additional processing required for non-documents
    if (!/documents/.test(resource)) {
      return baseDataProvider.update(resource, params);
    }

    console.log(resource, params);
    for (const key in params.data.values) {
      const value = params.data.values[key];
      // Ignore scalars
      if (typeof value != 'object') {
        continue;
      }

      let previousPath = '';
      if (params.previousData.hasOwnProperty(key)) {
        previousPath = params.previousData[key].path;
      }

      // Ignore already-uploaded files and unchanged paths
      if (!value.hasOwnProperty('rawFile') && previousPath === value.path) {
        console.log('Already uploaded and %s == %s', previousPath, value.path);
        continue;
      }

      let path = value.path;
      if (path === '') {
        // title is the file's basename
        path = value.title;
      }
      params.data.values[key] = uploadFile(path, value.rawFile).then((location) => ({
        url: location,
        path: path,
      }));
    }

    return baseDataProvider.update(resource, params);
  },
};

export default dataProvider;
