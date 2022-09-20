import { FormBuilder } from "@formio/react";
import { CommonInputProps, FieldTitle, useInput } from "react-admin";

export const FormBuilderInput = (props: CommonInputProps) => {
  const {
    defaultValue = '',
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
        form={field.value}
        options={{
          onchange: (newValue: any) => field.onChange(newValue),
        }}
      />
    </div>
  )
}
