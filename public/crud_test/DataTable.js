import Component from "../js-component-lib/lib/Component.js";
import { useTheme } from "../js-component-lib/lib/useTheme.js";
import TextField from "../js-component-lib/components/TextField.js";

class DataTable extends Component {
  handleDelete;
  handleUpdateRow;
  data;

  /**
   *
   * @param {{
   *  data: {id: string; number: number; boolean: boolean}[];
   *  handleDelete: (id: string) => {};
   *  handleUpdateRow: (id: string, number: number, boolean: boolean) => {};
   * }} props
   */
  constructor(props) {
    const { data, handleDelete, handleUpdateRow } = props;
    if (data == undefined || data == null) {
      throw new Error("data must be defined");
    }
    if (!Array.isArray(data)) {
      throw new Error("data must be an array");
    }

    super(document.createElement("table"));

    this.element.style.borderSpacing = "0px";

    this.data = data;
    this.handleDelete = handleDelete;
    this.handleUpdateRow = handleUpdateRow;

    this.initTheme();
    this.rebuild(data);
  }

  initTheme() {
    useTheme(
      () => {
        this.element.style.color = "white";
      },
      () => {
        this.element.style.color = "black";
      }
    );
  }

  /**
   *
   * @param {{id: string; number: number; boolean: boolean}[]} data
   */
  rebuild(data) {
    this.data = data;
    this.element.innerHTML = "";
    this.buildHeader();
    this.buildTable(data);
  }

  buildHeader() {
    const header = document.createElement("thead");
    const headerRow = document.createElement("tr");
    ["ID", "Number", "Boolean", "Edit", "Delete"].forEach((col) => {
      const cell = document.createElement("th");
      cell.innerHTML = col;
      cell.style.border = "1px solid";
      cell.style.padding = "5px";
      headerRow.appendChild(cell);
    });
    header.appendChild(headerRow);
    this.element.appendChild(header);
  }

  /**
   *
   * @param {HTMLTableCellElement} cell
   */
  applyCellStyles(cell) {
    cell.style.border = "1px solid";
    cell.style.padding = "5px";
    cell.style.textAlign = "center";
  }

  /**
   *
   * @param {{id: string; number: number; boolean: boolean}[]} data
   */
  buildTable(data) {
    const tbody = document.createElement("tbody");
    data.forEach((item) => {
      const tr = this.buildRow(item);
      tbody.appendChild(tr);
    });
    this.element.appendChild(tbody);
  }

  /**
   *
   * @param {{id: string; number: number; boolean: boolean}} item
   */
  buildRow(item) {
    const tr = document.createElement("tr");
    tr.id = item.id;
    const idCell = this.buildIdCell(item.id);
    const numberCell = this.buildNumberCell(item.number);
    const booleanCell = this.buildBooleanCell(item.boolean);
    const editCell = this.buildEditCell(item);
    const deleteCell = this.buildDeleteCell(item.id);
    tr.appendChild(idCell);
    tr.appendChild(numberCell);
    tr.appendChild(booleanCell);
    tr.appendChild(editCell);
    tr.appendChild(deleteCell);
    return tr;
  }

  /**
   *
   * @param {string} id
   * @returns
   */
  buildIdCell(id) {
    const idCell = document.createElement("td");
    idCell.innerHTML = id;
    this.applyCellStyles(idCell);
    return idCell;
  }

  /**
   *
   * @param {number} number
   * @returns
   */
  buildNumberCell(number) {
    const numberCell = document.createElement("td");
    numberCell.innerHTML = number.toString();
    this.applyCellStyles(numberCell);
    return numberCell;
  }

  /**
   *
   * @param {boolean} boolean
   * @returns
   */
  buildBooleanCell(boolean) {
    const booleanCell = document.createElement("td");
    booleanCell.innerHTML = boolean.toString();
    this.applyCellStyles(booleanCell);
    return booleanCell;
  }

  /**
   *
   * @param {{id: string; number: number; boolean: boolean;}} item
   * @returns
   */
  buildEditCell(item) {
    const editCell = document.createElement("td");
    const editButton = document.createElement("button");
    editButton.innerHTML = "Edit";
    editButton.onclick = () => this.editRow(item);
    editCell.appendChild(editButton);
    this.applyCellStyles(editCell);
    return editCell;
  }

  /**
   *
   * @param {string} id
   * @returns
   */
  buildDeleteCell(id) {
    const deleteCell = document.createElement("td");

    const deleteButton = document.createElement("button");
    deleteButton.innerHTML = "Delete";
    deleteButton.onclick = () => this.handleDelete(id);
    deleteCell.appendChild(deleteButton);

    this.applyCellStyles(deleteCell);

    return deleteCell;
  }

  /**
   *
   * @param {{id: string; number: number; boolean: boolean}} item
   */
  editRow(item) {
    const { id, number, boolean } = item;
    const tr = this.element.querySelector(`tr[id="${id}"]`);
    tr.innerHTML = "";

    const idCell = this.buildIdCell(id);
    const numberCell = this.buildNumberEditCell(id, number);
    const booleanCell = this.buildBooleanEditCell(id, boolean);
    const saveCell = this.buildSaveCell(id);
    const deleteCell = this.buildDeleteCell(id);

    tr.appendChild(idCell);
    tr.appendChild(numberCell);
    tr.appendChild(booleanCell);
    tr.appendChild(saveCell);
    tr.appendChild(deleteCell);
  }

  /**
   *
   * @param {string} id
   * @param {number} number
   * @returns
   */
  buildNumberEditCell(id, number) {
    const numberCell = document.createElement("td");
    const numberInput = new TextField({
      inputId: `number-${id}`,
      label: "Number",
      defaultValue: number.toString(),
    });
    numberCell.appendChild(numberInput.render());
    this.applyCellStyles(numberCell);
    return numberCell;
  }

  /**
   *
   * @param {string} id
   * @param {boolean} boolean
   * @returns
   */
  buildBooleanEditCell(id, boolean) {
    const booleanCell = document.createElement("td");
    const booleanInput = new TextField({
      inputId: `boolean-${id}`,
      label: "Boolean",
      defaultValue: boolean.toString(),
    });
    booleanCell.appendChild(booleanInput.render());
    this.applyCellStyles(booleanCell);
    return booleanCell;
  }

  /**
   *
   * @param {string} id
   * @returns
   */
  buildSaveCell(id) {
    const saveCell = document.createElement("td");
    const saveButton = document.createElement("button");
    saveButton.innerHTML = "Save";
    saveButton.onclick = () => {
      const numberValue = document.getElementById(`number-${id}`).value;
      let number = parseInt(numberValue);
      if (isNaN(number)) {
        number = 0;
      }
      const booleanValue = document.getElementById(`boolean-${id}`).value;
      let boolean;
      if (booleanValue == "true") {
        boolean = true;
      } else {
        boolean = false;
      }
      this.handleUpdateRow(id, number, boolean);
    };
    saveCell.appendChild(saveButton);
    this.applyCellStyles(saveCell);
    return saveCell;
  }
}

export default DataTable;
