import { render, screen } from "@testing-library/react";
import Home from "../app/page";

describe("Home", () => {
  it("renders correctly", () => {
    render(<Home />);
    // Check if some text from the default Next.js page exists
    const element = screen.getByText(/Get started/i);
    expect(element).toBeInTheDocument();
  });
});
