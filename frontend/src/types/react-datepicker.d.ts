declare module "react-datepicker" {
    import * as React from "react";
    export interface ReactDatePickerProps {
      selected?: Date | null;
      onChange?: (date: Date | null, event?: React.SyntheticEvent<any> | undefined) => void;
      showTimeSelect?: boolean;
      dateFormat?: string;
      className?: string;
      [key: string]: any;
    }
    export default class ReactDatePicker extends React.Component<ReactDatePickerProps> {}
  }
  