import React from 'react';
import { useInput } from 'ra-core';
import { Editor } from '@tinymce/tinymce-react';

export const TinyInput = (props) => {
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
    field,
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
    ...rest,
  });

  return (
    <Editor
      apiKey={process.env.REACT_APP_TINYMCE_KEY}
      initialValue={field.value}
      value={field.value}
      onEditorChange={field.onChange}
    />
  )
};
