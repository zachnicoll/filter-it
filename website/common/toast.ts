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

export const toastShow = (msg: string) => {
  return toast(msg, defaultToastOptions);
}

export const toastError = (msg: string, err: Error): ReactText => {
  console.error(err);
  return toast(msg, defaultToastOptions);
};

export const toastSuccess = (msg: string): ReactText => {
  toast.dismiss();
  return toast(msg, {
    ...defaultToastOptions,
    type: toast.TYPE.SUCCESS,
  });
};

export const toastUpdate = (toastId: string, opts: ToastOptions): void => {
  toast.dismiss();
  toast.update(toastId, {
    ...defaultToastOptions,
    ...opts,
  });
};

export const toastLoading = (msg: string): ReactText => {
  toast.dismiss();
  return toast.loading("Processing Image...", {...defaultToastOptions, autoClose: false});
};
