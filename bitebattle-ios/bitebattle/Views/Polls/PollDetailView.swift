import SwiftUI

struct PollDetailView: View {
    let poll: Poll

    @State private var results: [PollOptionResult] = []
    @State private var isLoading: Bool = false
    @State private var statusMessage: String?
    @State private var showAddOption: Bool = false

    var body: some View {
        AppBackground {
            VStack(spacing: 24) {
                headerView
                statusOrProgressView
                optionsListView
                Spacer()
                addButton
            }
            .padding()
            .navigationTitle("Poll Detail")
            .navigationBarTitleDisplayMode(.inline)
            .onAppear(perform: fetchResults)
        }
    }

    private var headerView: some View {
        Text(poll.name)
            .font(.largeTitle)
            .fontWeight(.bold)
            .foregroundColor(AppColors.textOnPrimary)
            .shadow(radius: 2)
            .padding(.top, 24)
    }

    private var statusOrProgressView: some View {
        Group {
            if isLoading {
                ProgressView()
                    .progressViewStyle(CircularProgressViewStyle(tint: AppColors.primary))
                    .padding()
            } else if let status = statusMessage {
                Text(status)
                    .foregroundColor(AppColors.textOnPrimary)
                    .font(.footnote)
                    .padding(.top, 8)
                    .shadow(radius: 1)
            } else if results.isEmpty {
                Text("No options yet.")
                    .foregroundColor(AppColors.textOnPrimary)
                    .font(.headline)
                    .padding()
                    .shadow(radius: 1)
            }
        }
    }

    private var optionsListView: some View {
        ScrollView {
            VStack(spacing: 18) {
                ForEach(results) { option in
                    AppTile {
                        VStack(alignment: .leading, spacing: 8) {
                            HStack {
                                Text(option.option_name)
                                    .font(.headline)
                                    .foregroundColor(AppColors.secondary)
                                Spacer()
                                Text("Votes: \(option.vote_count)")
                                    .font(.subheadline)
                                    .foregroundColor(AppColors.textOnPrimary)
                                    .padding(.horizontal, 8)
                                    .padding(.vertical, 2)
                                    .background(AppColors.secondary.opacity(0.7))
                                    .cornerRadius(8)
                            }
                            if !option.voter_ids.isEmpty {
                                Text("Voters: \(option.voter_ids.joined(separator: ", "))")
                                    .font(.caption)
                                    .foregroundColor(AppColors.textOnPrimary)
                                    .padding(.horizontal, 6)
                                    .padding(.vertical, 2)
                                    .background(AppColors.primary.opacity(0.5))
                                    .cornerRadius(6)
                            }
                        }
                    }
                    .transition(.move(edge: .bottom).combined(with: .opacity))
                }
            }
            .padding(.horizontal, 12)
            .animation(.easeInOut(duration: 0.18), value: results)
        }
    }

    private var addButton: some View {
        AppButton(
            title: "Add Option",
            icon: "plus.circle.fill",
            background: AppColors.primary,
            foreground: AppColors.textOnPrimary,
            isLoading: false,
            isDisabled: false
        ) {
            showAddOption = true
        }
        .padding(.horizontal, 12)
        .padding(.bottom, 24)
        .sheet(isPresented: $showAddOption, onDismiss: fetchResults) {
            PollOptionView(poll: poll)
        }
    }

    private func fetchResults() {
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

// Conform PollOptionResult to Equatable for animation
extension PollOptionResult: Equatable {
    static func == (lhs: PollOptionResult, rhs: PollOptionResult) -> Bool {
        lhs.option_id == rhs.option_id &&
        lhs.option_name == rhs.option_name &&
        lhs.vote_count == rhs.vote_count &&
        lhs.voter_ids == rhs.voter_ids
    }
}
