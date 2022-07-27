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
    SelectInput,
    FormDataConsumer,
    DateInput,
    DateTimeInput,
    ReferenceInput
} from 'react-admin';

const fieldChoices = [
    { id: 'date',          name: 'Date' },
    { id: 'datetime',      name: 'Date & Time' },
    { id: 'email',         name: 'Email' },
    { id: 'select',        name: 'Multi-Select' },
    { id: 'number',        name: 'Number' },
    // { id: 'select-class',  name: 'Select (Class)' }, // Need data_source_* fields on Field struct
    { id: 'select-static', name: 'Select (Static)' },
    { id: 'text',          name: 'Text' },
    { id: 'textarea',      name: 'Textarea' },
    { id: 'time',          name: 'Time' },
    { id: 'tinymce',       name: 'TinyMCE' },
    { id: 'upload',        name: 'Upload' },
];

export const ClassCreate = (props) => {
    const ensureFields = data => ({
        ...data,
        fields: [],
    });

    return (
        <Create {...props} transform={ensureFields}>
            <SimpleForm>
                <TextInput source="name" validate={[required()]} fullWidth />
                <TextInput source="table_labels" />
                <TextInput source="table_fields" />
            </SimpleForm>
        </Create>
    );
};

export const ClassEdit = (props) => (
    <Edit {...props}>
        <SimpleForm>
            <TextInput source="name" validate={[required()]} fullWidth />
            <TextInput source="table_labels" />
            <TextInput source="table_fields" />
            <ArrayInput source="fields">
                <SimpleFormIterator >
                    <TextInput source="label" />
                    <TextInput source="name" />
                    <SelectInput source="type" choices={fieldChoices} />
                    <FormDataConsumer>
                        {({
                            formData,
                            scopedFormData,
                            getSource,
                            ...rest
                        }) => {
                            switch (scopedFormData.type) {
                            case 'date':
                                return (
                                    <>
                                        <DateInput source={getSource('min')} {...rest} />
                                        <DateInput source={getSource('max')} {...rest} />
                                        <TextInput source={getSource('step')} label="Step (days)" {...rest} />
                                        <TextInput source={getSource('format')} label="Format (Jan 2, 2006 3:04pm)" {...rest} />
                                    </>
                                );
                            case 'datetime':
                                return (
                                    <>
                                        <DateTimeInput source={getSource('min')} {...rest} />
                                        <DateTimeInput source={getSource('max')} {...rest} />
                                        <TextInput source={getSource('step')} label="Step (days)" {...rest} />
                                        <TextInput source={getSource('format')} label="Format (Jan 2, 2006 3:04pm)" {...rest} />
                                    </>
                                );
                            case 'time':
                                return (
                                    <>
                                        <TextInput source={getSource('format')} label="Format (3:04pm)" {...rest} />
                                    </>
                                );
                            case 'number':
                                return (
                                    <>
                                        <TextInput source={getSource('min')} {...rest} />
                                        <TextInput source={getSource('max')} {...rest} />
                                        <TextInput source={getSource('step')} {...rest} />
                                    </>
                                );
                            case 'select-static':
                                return (
                                    <>
                                        <TextInput source={getSource('options')} label="Options (one per line, key | value or just value" multiline {...rest} />
                                    </>
                                );
                            case 'select-class':
                                return (
                                    <>
                                        <ReferenceInput source={getSource('data_source_id')} reference="classes">
                                            <SelectInput optionText="name" />
                                        </ReferenceInput>
                                        <TextInput source={getSource('data_source_value')} />
                                        <TextInput source={getSource('data_source_label')} />
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

export const ClassList = () => (
    <List>
        <Datagrid rowClick="edit">
            <TextField source="name" />
            <DateField source="created" />
            <DateField source="updated" />
            <EditButton />
        </Datagrid>
    </List>
);
