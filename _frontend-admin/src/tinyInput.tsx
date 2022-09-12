import { useInput, FieldTitle } from 'ra-core';
import { Editor } from '@tinymce/tinymce-react';
import { CommonInputProps } from 'react-admin';

// This implementation may change if there are complaints about performance
// later on. There is a way to update the content when "dirty" instead of every
// keystroke as the current implementation does.
// See: https://www.tiny.cloud/docs/tinymce/6/react-ref/#using-the-tinymce-react-component-as-a-uncontrolled-component
export const TinyInput = (props: CommonInputProps) => {
  const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;

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

  const init = {
    skin: (prefersDark ? 'oxide-dark' : ''),
    content_css: (prefersDark ? 'dark' : ''),
    width: '100%',
  };

  return (
    <div style={{ marginBottom: '1em', width: '100%' }}>
      <FieldTitle label={label} source={source} resource={resource} />
      <Editor
        apiKey={process.env.REACT_APP_TINYMCE_KEY}
        init={init}
        plugins={['code']}
        value={field.value}
        onEditorChange={(newValue) => field.onChange(newValue)}
      />
    </div>
  )
};
