import simpleRestProvider from 'ra-data-simple-rest';

const API_URL = process.env.REACT_APP_API_URL;
const baseDataProvider = simpleRestProvider(API_URL);

const dataProvider = {
  ...baseDataProvider,
  update: (resource, params) => {
    // No additional processing required for non-documents
    if (!/documents/.test(resource)) {
      return baseDataProvider.update(resource, params);
    }

    console.log(resource, params);
    let files = [];
    for (const key in params.data.values) {
      const value = params.data.values[key];
      // Ignore scalars
      if (typeof value != 'object') {
        continue;
      }

      let previousPath = '';
      if (params.data.previousData.hasOwnProperty(key)) {
        previousPath = params.data.previousData[key].path;
      }

      // Ignore already-uploaded files and unchanged paths
      if (!value.hasOwnProperty('rawFile') && previousPath == value.path) {
        console.log('Already uploaded and %s == %s', previousPath, value.path);
        continue;
      }
      console.log(value);

      return baseDataProvider.update(resource, params);
    }
  },
};

export default dataProvider;
