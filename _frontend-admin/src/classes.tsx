import React from 'react';
import {
  ArrayInput,
  BooleanInput,
  Create,
  CreateProps,
  Datagrid,
  DateField,
  DateInput,
  DateTimeInput,
  Edit,
  EditButton,
  EditProps,
  FormDataConsumer,
  FormDataConsumerRenderParams,
  List,
  ListProps,
  NumberInput,
  ReferenceInput,
  SelectInput,
  SimpleForm,
  SimpleFormIterator,
  TextField,
  TextInput,
  TransformData,
  required,
  useRedirect,
} from 'react-admin';
import './App.css';

interface UpdateProps {
  update: React.Dispatch<React.SetStateAction<number>>;
}

interface CreateUpdateProps extends CreateProps, UpdateProps {};

interface EditUpdateProps extends EditProps, UpdateProps {};

const fieldChoices = [
  { id: 'date',               name: 'Date' },
  { id: 'datetime',           name: 'Date & Time' },
  { id: 'tiny',               name: 'Editor' },
  { id: 'email',              name: 'Email' },
  { id: 'number',             name: 'Number' },
  { id: 'select-class',       name: 'Select (Class)' },
  { id: 'multi-class',        name: 'Mutli-Select (Class)' },
  { id: 'multi-select-label', name: 'Multi-Select w/ Label' },
  { id: 'select-static',      name: 'Select (Static)' },
  { id: 'text',               name: 'Text' },
  { id: 'textarea',           name: 'Textarea' },
  { id: 'time',               name: 'Time' },
  { id: 'any-upload',         name: 'Upload (Any)' },
  { id: 'image-upload',       name: 'Upload (Image)' },
];

export const ClassCreate = (props: CreateUpdateProps) => {
  const { update, ...rest } = props;
  const redirect = useRedirect();
  const onSuccess = () => {
    update((new Date()).getTime());
    redirect('list', 'classes');
  };

  const ensureFields = (data: TransformData) => ({
    ...data,
    fields: [],
  });

  return (
    <Create {...rest} mutationOptions={{ onSuccess }} transform={ensureFields}>
      <SimpleForm>
        <TextInput source="name" validate={[required()]} fullWidth />
      </SimpleForm>
    </Create>
  );
};

export const ClassEdit = (props: EditUpdateProps) => {
  const { update, ...rest } = props;
  const redirect = useRedirect();
  const onSuccess = () => {
    update((new Date()).getTime());
    redirect('list', 'classes');
  };

  return (
    <Edit {...rest} mutationOptions={{ onSuccess }} mutationMode="pessimistic">
      <SimpleForm>
        <TextInput source="name" validate={[required()]} fullWidth />
        <ReferenceInput source="parent_id" reference="classes" >
          <SelectInput optionText="name" fullWidth />
        </ReferenceInput>
        <ArrayInput source="fields">
          <SimpleFormIterator className="field-row">
            <TextInput source="label" />
            <TextInput source="name" />
            <BooleanInput source="sort" />
            <NumberInput source="column" />
            <SelectInput source="type" choices={fieldChoices} defaultValue="text" />
            <FormDataConsumer>
              {({
                scopedFormData,
                getSource,
              }: FormDataConsumerRenderParams) => {
                const getSrc = getSource || ((s: string) => s);
                switch (scopedFormData.type) {
                case 'date':
                  return (
                    <>
                      <DateInput source={getSrc('min')} />
                      <DateInput source={getSrc('max')} />
                      <TextInput source={getSrc('step')} label="Step (days)" />
                      <TextInput source={getSrc('format')} label="Format (Jan 2, 2006 3:04pm)" />
                    </>
                  );
                case 'datetime':
                  return (
                    <>
                      <DateTimeInput source={getSrc('min')} />
                      <DateTimeInput source={getSrc('max')} />
                      <TextInput source={getSrc('step')} label="Step (days)" />
                      <TextInput source={getSrc('format')} label="Format (Jan 2, 2006 3:04pm)" />
                    </>
                  );
                case 'time':
                  return (
                    <>
                      <TextInput source={getSrc('format')} label="Format (3:04pm)" />
                    </>
                  );
                case 'number':
                  return (
                    <>
                      <TextInput source={getSrc('min')} />
                      <TextInput source={getSrc('max')} />
                      <TextInput source={getSrc('step')} />
                    </>
                  );
                case 'select-static':
                  return (
                    <>
                      <TextInput source={getSrc('options')} label="Options (one per line, key | value or just value" multiline />
                    </>
                  );
                case 'select-class':
                case 'multi-class':
                case 'multi-select-label':
                  return (
                    <>
                      <ReferenceInput source={getSrc('class_id')} reference="classes">
                        <SelectInput optionText="name" />
                      </ReferenceInput>
                      <TextInput source={getSrc('field')} />
                    </>
                  );
                default:
                  return null;
                }
              }}
            </FormDataConsumer>
          </SimpleFormIterator>
        </ArrayInput>
      </SimpleForm>
    </Edit>
  );
};

export const ClassList = (props: ListProps) => (
  <List {...props}>
    <Datagrid>
      <TextField source="name" />
      <DateField source="created" />
      <DateField source="updated" />
      <EditButton />
    </Datagrid>
  </List>
);
