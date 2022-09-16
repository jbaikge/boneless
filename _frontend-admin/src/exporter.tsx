export const jsonExporter = (baseName: string) => (objects: Array<any>) => {
  const blob = new Blob([JSON.stringify(objects, null, 2)], { type: 'application/json' });
  const link = document.createElement('a');
  link.href = window.URL.createObjectURL(blob);
  link.download = `${baseName}.json`;
  link.click();
};
