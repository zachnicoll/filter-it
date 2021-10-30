import { ReactText } from "react";
import { toast, ToastOptions } from "react-toastify";
import theme from "styles/colors";

const defaultToastOptions: ToastOptions = {
  position: "bottom-right",
  autoClose: 5000,
  hideProgressBar: true,
  closeOnClick: true,
  pauseOnHover: true,
  draggable: true,
  progress: undefined,
  style: {
    borderColor: theme.accent,
    borderWidth: 1,
    borderStyle: "solid",
  },
};

export const toastError = (msg: string, err: Error): ReactText => {
  console.error(err);
  return toast(msg, defaultToastOptions);
};

export const toastUpdate = (msg: string): void => {
  toast.dismiss();
  toast.update(msg, {
    ...defaultToastOptions,
    type: toast.TYPE.SUCCESS,
  });
};

export const toastLoading = (msg: string): ReactText => {
  toast.dismiss();
  return toast.loading("Processing Image...", defaultToastOptions);
};
