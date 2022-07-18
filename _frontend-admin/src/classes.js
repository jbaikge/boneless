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
    Edit,
    ArrayInput,
    SimpleFormIterator,
    SelectInput
} from 'react-admin';

const fieldChoices = [
    { id: 'date', name: 'Date' },
    { id: 'text', name: 'Text' },
    { id: 'time', name: 'Time' },
];

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
            <ArrayInput source="fields">
                <SimpleFormIterator>
                    <TextInput source="label" />
                    <TextInput source="name" />
                    <SelectInput source="type" choices={fieldChoices} />
                </SimpleFormIterator>
            </ArrayInput>
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