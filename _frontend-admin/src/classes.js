import {
    Datagrid,
    DateField,
    List,
    TextField,
    Create,
    SimpleForm,
    TextInput,
    required,
    EditButton,
    Edit
} from 'react-admin';

export const ClassCreate = () => (
    <Create>
        <SimpleForm>
            <TextInput source="name" validate={[required()]} fullWidth />
            <TextInput source="slug" validate={[required()]} fullWidth />
            <TextInput source="table_labels" />
            <TextInput source="table_fields" />
        </SimpleForm>
    </Create>
)

export const ClassEdit = (props) => (
    <Edit {...props}>
        <SimpleForm>
            <TextInput source="name" validate={[required()]} fullWidth />
            <TextInput source="slug" validate={[required()]} fullWidth />
            <TextInput source="table_labels" />
            <TextInput source="table_fields" />
        </SimpleForm>
    </Edit>
)

export const ClassList = () => (
    <List>
        <Datagrid rowClick="edit">
            <TextField source="id" />
            <TextField source="name" />
            <TextField source="slug" />
            <DateField source="created" />
            <DateField source="updated" />
            <EditButton />
        </Datagrid>
    </List>
);