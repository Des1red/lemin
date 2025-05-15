import json
import argparse
import matplotlib.pyplot as plt
import networkx as nx
from collections import defaultdict

def main():
    parser = argparse.ArgumentParser(description="Visualize lem-in simulation")
    parser.add_argument('--input', '-i', required=True, help='Simulation JSON file')
    args = parser.parse_args()

    # Load simulation data
    with open(args.input) as f:
        data = json.load(f)
    rooms = data['rooms']
    moves = data['moves']

    # Build graph: nodes with positions; edges inferred from move history
    G = nx.Graph()
    pos = {}
    for r in rooms:
        G.add_node(r['name'])
        pos[r['name']] = (r['x'], r['y'])
    # infer edges from move history
    edges = set()
    for m in moves:
        frm = m['from']
        to = m['to']
        if frm and to:
            edges.add((frm, to))
    G.add_edges_from(edges)

    # Organize moves by turn
    moves_by_turn = defaultdict(list)
    for m in moves:
        moves_by_turn[m['turn']].append(m)

    # Track current positions of ants
    ant_positions = {}

    # Setup plot
    plt.ion()
    fig, ax = plt.subplots()
    ax.set_title("Lem-in Ant Simulation")

    # Iterate through turns
    for turn in sorted(moves_by_turn.keys()):
        # Update ant positions for this turn
        for m in moves_by_turn[turn]:
            ant_positions[m['ant']] = m['to']

        ax.clear()
        # Draw graph structure
        nx.draw(G, pos=pos, ax=ax, with_labels=True, node_size=500, font_size=8)
        # Draw ants as red text at their node positions
        for ant, room in ant_positions.items():
            x, y = pos[room]
            ax.text(x, y, f"{ant}", fontsize=10, fontweight="bold",
                    ha='center', va='center', color='red')

        ax.set_title(f"Turn {turn}")
        plt.pause(0.5)

    plt.ioff()
    plt.show()

if __name__ == '__main__':
    main()
