import React from "react";
import * as Styles from "../UploadBox.styles";

interface DragNDropProps {
  getRootProps: () => any;
  getInputProps: () => any;
  openDialog: () => void;
  dragActive: boolean;
}

export const DragNDrop: React.FC<DragNDropProps> = ({
  getRootProps,
  getInputProps,
  openDialog,
  dragActive,
}) => {
  return (
    <Styles.Container {...getRootProps()} active={dragActive}>
      <input {...getInputProps()} />
      <p>Drag &apos;n&apos; drop some files here, or click to select files</p>
      <br />
      <button type="button" onClick={openDialog}>
        Open File Dialog
      </button>
    </Styles.Container>
  );
};
