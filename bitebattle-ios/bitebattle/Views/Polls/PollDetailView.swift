import SwiftUI

struct PollDetailView: View {
    let poll: PollsView.Poll

    @State private var results: [PollOptionResult] = []
    @State private var isLoading: Bool = false
    @State private var statusMessage: String?
    @State private var showAddOption: Bool = false

    struct PollOptionResult: Identifiable, Decodable {
        let option_id: String
        let option_name: String
        let vote_count: Int
        let voter_ids: [String]

        var id: String { option_id }

        private enum CodingKeys: String, CodingKey {
            case option_id, option_name, vote_count, voter_ids
        }

        init(from decoder: Decoder) throws {
            let container = try decoder.container(keyedBy: CodingKeys.self)
            option_id = try container.decode(String.self, forKey: .option_id)
            option_name = try container.decode(String.self, forKey: .option_name)
            vote_count = try container.decode(Int.self, forKey: .vote_count)
            voter_ids = try container.decodeIfPresent([String].self, forKey: .voter_ids) ?? []
        }
    }

    var body: some View {
        ZStack {
            LinearGradient(
                gradient: Gradient(colors: [Color.orange.opacity(0.7), Color.pink.opacity(0.7)]),
                startPoint: .topLeading,
                endPoint: .bottomTrailing
            )
            .ignoresSafeArea()

            VStack(spacing: 24) {
                Text(poll.name)
                    .font(.largeTitle)
                    .fontWeight(.bold)
                    .foregroundColor(.white)
                    .shadow(radius: 4)
                    .padding(.top, 24)

                if isLoading {
                    ProgressView()
                        .progressViewStyle(CircularProgressViewStyle(tint: .white))
                        .padding()
                } else if let status = statusMessage {
                    Text(status)
                        .foregroundColor(.white)
                        .font(.footnote)
                        .padding(.top, 8)
                        .shadow(radius: 1)
                } else if results.isEmpty {
                    Text("No options yet.")
                        .foregroundColor(.white)
                        .font(.headline)
                        .padding()
                        .shadow(radius: 1)
                }

                ScrollView {
                    VStack(spacing: 18) {
                        ForEach(results) { option in
                            VStack(alignment: .leading, spacing: 8) {
                                HStack {
                                    Text(option.option_name)
                                        .font(.headline)
                                        .foregroundColor(.orange)
                                    Spacer()
                                    Text("Votes: \(option.vote_count)")
                                        .font(.subheadline)
                                        .foregroundColor(.white)
                                        .padding(.horizontal, 8)
                                        .padding(.vertical, 2)
                                        .background(Color.orange.opacity(0.7))
                                        .cornerRadius(8)
                                }
                                if !option.voter_ids.isEmpty {
                                    Text("Voters: \(option.voter_ids.joined(separator: ", "))")
                                        .font(.caption)
                                        .foregroundColor(.white)
                                        .padding(.horizontal, 6)
                                        .padding(.vertical, 2)
                                        .background(Color.pink.opacity(0.5))
                                        .cornerRadius(6)
                                }
                            }
                            .padding()
                            .background(Color.white.opacity(0.15))
                            .cornerRadius(14)
                            .shadow(radius: 2)
                        }
                    }
                    .padding(.horizontal, 12)
                }

                Spacer()

                Button(action: {
                    showAddOption = true
                }) {
                    HStack {
                        Image(systemName: "plus.circle.fill")
                            .foregroundColor(.white)
                        Text("Add Option")
                            .fontWeight(.semibold)
                            .foregroundColor(.white)
                    }
                    .frame(maxWidth: .infinity, minHeight: 48)
                    .padding()
                    .background(Color.pink.opacity(0.8))
                    .cornerRadius(14)
                    .shadow(radius: 2)
                }
                .padding(.horizontal, 12)
                .padding(.bottom, 24)
                .sheet(isPresented: $showAddOption, onDismiss: fetchResults) {
                    PollOptionView(poll: poll)
                }
            }
            .padding(.top)
        }
        .navigationTitle("Poll Details")
        .navigationBarTitleDisplayMode(.inline)
        .onAppear(perform: fetchResults)
    }

    func fetchResults() {
        guard let token = UserDefaults.standard.string(forKey: "authToken"),
              !token.isEmpty,
              let url = URL(string: "http://localhost:8080/api/polls/\(poll.id)/results") else {
            statusMessage = "Not logged in."
            return
        }

        var request = URLRequest(url: url)
        request.httpMethod = "GET"
        request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")

        isLoading = true
        statusMessage = nil
        results = []

        URLSession.shared.dataTask(with: request) { data, response, error in
            DispatchQueue.main.async {
                isLoading = false
                if let error = error {
                    statusMessage = "Failed: \(error.localizedDescription)"
                    return
                }
                guard let data = data else {
                    statusMessage = "No data returned."
                    return
                }
                do {
                    let decoded = try JSONDecoder().decode([PollOptionResult].self, from: data)
                    self.results = decoded
                } catch {
                    statusMessage = "Failed to load results"
                }
            }
        }.resume()
    }
}
