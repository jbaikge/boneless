import simpleRestProvider from 'ra-data-simple-rest';
import {
  CreateParams,
  UpdateParams,
  fetchUtils
} from 'ra-core';

interface FileProps {
  key: string;
  path: string;
  file: File;
  location?: string;
};

interface FileReaderResult {
  url: string;
  method: string;
  body: string | ArrayBuffer | null;
}

const API_URL: string = process.env.REACT_APP_API_URL || '';
const baseDataProvider = simpleRestProvider(API_URL);
const documentRE = /documents/;

const uploadFile = (fileInfo: FileProps) =>
  fetchUtils.fetchJson(`${API_URL}/files/url`, {
    method: 'POST',
    body: JSON.stringify({
      key: fileInfo.path,
      expires: '5m',
      content_type: fileInfo.file.type,
    }),
  }).then(({ json }): Promise<FileReaderResult> => new Promise((resolve, reject) => {
    fileInfo.location = json.location;
    const reader = new FileReader();
    reader.addEventListener('load', () => resolve({
      url: json.url,
      method: json.method,
      body: reader.result
    }));
    reader.addEventListener('error', reject);
    reader.readAsArrayBuffer(fileInfo.file);
  })).then(({ url, method, body }) => fetchUtils.fetchJson(url, {
    method: method,
    headers: new Headers({ 'Content-Type': fileInfo.file.type }),
    body: body,
  })).then(() => fileInfo);

const dataProvider = {
  ...baseDataProvider,
  create: (resource: string, params: CreateParams) => {
    // No additional processing required for non-documents
    if (!documentRE.test(resource)) {
      return baseDataProvider.create(resource, params);
    }

    let files = [];
    for (const key in params.data.values) {
      const value = params.data.values[key];

      // Ignore scalars
      if (typeof value != 'object') {
        continue;
      }

      // Ignore objects without files
      if (!value.hasOwnProperty('rawFile')) {
        continue;
      }

      // Use the original filename if no path provided
      let path = value.path;
      if (path === '') {
        path = value.title;
      }

      files.push({
        key: key,
        path: path,
        file: value.rawFile,
      });
    }

    return Promise.all(files.map(uploadFile))
    .then((infos) => infos.map((info) => {
      return params.data.values[info.key] = {
        path: info.path,
        url: info.location,
      }
    }))
    .then(() => baseDataProvider.create(resource, params));
  },
  update: (resource: string, params: UpdateParams) => {
    // No additional processing required for non-documents
    if (!documentRE.test(resource)) {
      return baseDataProvider.update(resource, params);
    }

    let files: Array<FileProps> = [];
    for (const key in params.data.values) {
      const value = params.data.values[key];

      // Ignore scalars
      if (typeof value != 'object') {
        continue;
      }

      // Ignore objects without files
      if (!value.hasOwnProperty('rawFile')) {
        continue;
      }

      let previousPath = '';
      if (params.previousData.hasOwnProperty(key)) {
        previousPath = params.previousData[key].path;
      }

      // Ignore already-uploaded files and unchanged paths
      if (!value.hasOwnProperty('rawFile') && previousPath === value.path) {
        continue;
      }

      // Use the original filename if no path provided
      let path = value.path;
      if (path === '') {
        path = value.title;
      }

      files.push({
        key: key,
        path: path,
        file: value.rawFile,
      });
    }

    return Promise.all(files.map(uploadFile))
      .then((infos) => infos.map((info) => {
        return params.data.values[info.key] = {
          path: info.path,
          url: info.location,
        }
      }))
      .then(() => baseDataProvider.update(resource, params));
  },
};

export default dataProvider;
