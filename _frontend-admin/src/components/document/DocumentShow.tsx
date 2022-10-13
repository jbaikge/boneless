import {
  ArrayField,
  Datagrid,
  ImageField,
  Loading,
  ReferenceField,
  RichTextField,
  Show,
  ShowProps,
  SimpleShowLayout,
  TextField,
  useGetOne,
  useResourceContext,
} from 'react-admin';
import { reResource } from './Constants';
import { FieldProps } from '../field/Props';

const DocumentShow = (props: ShowProps) => {
  const resourceContext = useResourceContext();
  const [[ , resource, id ]] = [...resourceContext.matchAll(reResource)];
  const { data, isLoading } = useGetOne(resource, { id });

  if (isLoading) {
    return <Loading />;
  }

  let parentField = null;
  if (data.parent_id !== '') {
    parentField = (
      <ReferenceField reference={`classes/${data.parent_id}/documents`} source="parent_id" label="Parent" sortable={false}>
        <TextField source="values.title" />
      </ReferenceField>
    );
  }

  return (
    <Show {...props}>
      <SimpleShowLayout>
        {parentField}
        {data.fields.map((field: FieldProps) => {
          const source = `values.${field.name}`;
          const label = `${field.label} (.Document.Values.${field.name})`
          switch (field.type) {
            case 'image-upload':
              return <ImageField source={`${source}.url`} label={label} />
            case 'multi-class':
              return <ArrayField source={source} label={false}>
                <Datagrid bulkActionButtons={false}>
                  <ReferenceField reference={`classes/${field.class_id}/documents`} source="id" label={label}>
                    <TextField source={`values.${field.field}`} />
                  </ReferenceField>
                </Datagrid>
              </ArrayField>
            case 'tiny':
              return <RichTextField source={source} label={label} />
            default:
              return <TextField source={source} label={label} />
          }
        })}
      </SimpleShowLayout>
    </Show>
  )
}

export default DocumentShow;
