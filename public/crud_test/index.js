import Body from "../js-component-lib/components/Body.js";
import DataTable from "./DataTable.js";
import Div from "../js-component-lib/components/Div.js";
import AddButton from "./AddButton.js";

window.onload = () => {
  document.title = "CRUD Test";

  /**
   *
   * @param {string} id
   */
  const handleDeleteRow = async (id) => {
    try {
      await fetch(`/test/${id}/delete`, {
        method: "DELETE",
      });
      dataTable.rebuild(dataTable.data.filter((i) => i.id != id));
    } catch (e) {
      console.error(e);
    }
  };

  const dataTable = new DataTable({ data: [], handleDelete: handleDeleteRow });

  const getData = async () => {
    try {
      const response = await fetch("/test/read", {
        method: "GET",
      });
      const json = await response.json();
      dataTable.rebuild(json);
    } catch (e) {
      console.error(e);
    }
  };

  const handleAddRow = async () => {
    try {
      const response = await fetch("/test/create", {
        method: "POST",
        body: JSON.stringify({
          number: 1,
          boolean: true,
        }),
      });
      const json = await response.json();
      dataTable.rebuild([...dataTable.data, json]);
    } catch (e) {
      console.error(e);
    }
  };

  getData();

  new Body({
    children: [
      new Div({
        styles: {
          width: "100vw",
          height: "max-content",
          boxSizing: "border-box",
          margin: "0px",
          padding: "10px",
          display: "flex",
          gap: "1rem",
        },
        children: [
          dataTable,
          new AddButton({
            handleAdd: handleAddRow,
          }),
        ],
      }),
    ],
  });
};
