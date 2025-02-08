import Component from "../js-component-lib/lib/Component.js";

class AddButton extends Component {
  /**
   *
   * @param {{handleAdd: () => {}}} props
   */
  constructor(props) {
    const { handleAdd } = props;

    super(document.createElement("button"));

    this.element.innerHTML = "Add Row";
    this.element.onclick = handleAdd;

    this.element.style.height = "30px";
    this.element.style.width = "100px";
  }
}

export default AddButton;
