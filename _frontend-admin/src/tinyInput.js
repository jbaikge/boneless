import React from 'react';
import { useInput } from 'ra-core';
import { Editor } from '@tinymce/tinymce-react';

// This implementation may change if there are complaints about performance
// later on. There is a way to update the content when "dirty" instead of every
// keystroke as the current implementation does.
// See: https://www.tiny.cloud/docs/tinymce/6/react-ref/#using-the-tinymce-react-component-as-a-uncontrolled-component
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
      init={{
        width: '100%',
      }}
      plugins={['code']}
      value={field.value}
      onEditorChange={(newValue) => field.onChange(newValue)}
    />
  )
};
