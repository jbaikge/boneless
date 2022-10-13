import {
  Datagrid,
  DateField,
  EditButton,
  ImageField,
  List,
  ListProps,
  Loading,
  ReferenceField,
  ReferenceInput,
  SearchInput,
  SelectInput,
  ShowButton,
  TextField,
  useResourceContext,
  useGetOne,
} from 'react-admin';
import { reResource } from './Constants';
import { FieldProps } from '../field/Props';
import GlobalPagination from '../GlobalPagination';

const DocumentList = (props: ListProps) => {
  const resourceContext = useResourceContext();
  const [[ , resource, id ]] = [...resourceContext.matchAll(reResource)];
  const { data, isLoading } = useGetOne(resource, { id });

  if (isLoading) {
    return <Loading />;
  }

  const filters = [
    <SearchInput source="q" alwaysOn />
  ];

  let parentField = null;
  if (data.parent_id !== '') {
    parentField = (
      <ReferenceField reference={`classes/${data.parent_id}/documents`} source="parent_id" label="Parent" sortable={false}>
        <TextField source="values.title" />
      </ReferenceField>
    );
    filters.push(
      <ReferenceInput reference={`classes/${data.parent_id}/documents`} source="parent_id" label="Parent" alwaysOn>
        <SelectInput optionText="values.title" />
      </ReferenceInput>
    );
  }

  return (
    <List {...props} pagination={<GlobalPagination />} filters={filters}>
      <Datagrid sx={{
        '& td:last-child': { width: '5em' },
        '& td:nth-last-of-type(2)': { width: '5em' },
      }}>
        {parentField}
        {data.fields.filter((field: FieldProps) => field.column > 0).sort((a: FieldProps, b: FieldProps) => a.column - b.column).map((field: FieldProps) => {
          const source = `values.${field.name}`;
          switch (field.type) {
            case 'date':
              return <DateField key={field.label} source={source} label={field.label} sortable={field.sort} />
            case 'datetime':
              return <DateField key={field.label} source={source} label={field.label} sortable={field.sort} showTime />
            case 'image-upload':
              return <ImageField key={field.label} source={`${source}.url`} label={field.label} sortable={field.sort} />
            default:
              return <TextField key={field.label} source={source} label={field.label} sortable={field.sort} />
          }
        })}
        <ShowButton />
        <EditButton />
      </Datagrid>
    </List>
  );
};

export default DocumentList;
