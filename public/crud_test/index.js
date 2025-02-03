import Body from "../js-component-lib/components/Body.js";
import DataTable from "./DataTable.js";
import Div from "../js-component-lib/components/Div.js";

window.onload = () => {
  document.title = "CRUD Test";

  new Body({
    children: [
      new Div({
        styles: {
          width: "100%",
          height: "100%",
          boxSizing: "border-box",
          margin: "0px",
          padding: "10px",
        },
        children: [
          new DataTable({
            data: [],
          }),
        ],
      }),
    ],
  });
};
