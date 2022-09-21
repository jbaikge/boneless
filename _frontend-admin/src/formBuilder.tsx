import { FormBuilder } from "@formio/react";
import { CommonInputProps, FieldTitle, useInput } from "react-admin";

export const FormBuilderInput = (props: CommonInputProps) => {
  const {
    defaultValue = {display: "form"},
    format,
    parse,
    resource,
    source,
    validate,
    onBlur,
    onChange,
    label,
    helperText,
    ...rest
  } = props;

  const {
    field
  } = useInput({
    defaultValue,
    format,
    parse,
    resource,
    source,
    type: 'text',
    validate,
    onBlur,
    onChange,
    ...rest
  });

  return (
    <div style={{ marginBottom: '1em', width: '100%' }}>
      <FieldTitle label={label} source={source} resource={resource} />
      <FormBuilder
        onChange={(schema: any) => field.onChange(schema)}
        onUpdateComponent={(component: any, something: any, schema: any) => console.log('*** onUpdateComponent ***', 'component', component, 'something', something, 'schema', schema)}
        onSaveComponent={(component: any, something: any, schema: any) => console.log('*** onSaveComponent ***', 'component', component, 'something', something, 'schema', schema)}
        onEditComponent={(component: any, something: any, schema: any) => console.log('*** onEditComponent ***', 'component', component, 'something', something, 'schema', schema)}
        onDeleteComponent={(component: any, something: any, schema: any) => console.log('*** onDeleteComponent ***', 'component', component, 'something', something, 'schema', schema)}
        onCancelComponent={(component: any, something: any, schema: any) => console.log('*** onCancelComponent ***', 'component', component, 'something', something, 'schema', schema)}
        form={field.value}
        options={{
          onchange: (newValue: any) => console.log('onchange', newValue),
          onChange: (one: any, two: any, three: any) => console.log('onChange', one, two, three),
        }}
      />
    </div>
  )
}
