import Foundation
import SwiftUI


struct PollsView: View {
    @State private var polls: [Poll] = []
    @State private var isLoading: Bool = false
    @State private var statusMessage: String?
    @State private var showAddPoll: Bool = false
    @State private var newPollName: String = ""
    @State private var isCreatingPoll: Bool = false
    @State private var showJoinPoll: Bool = false
    @State private var inviteCode: String = ""
    @State private var isJoiningPoll: Bool = false

    var body: some View {
        AppBackground {
            VStack(spacing: 20) {
                TitleText(title: "Your Polls")
                StatusOrLoadingView(
                    isLoading: isLoading,
                    statusMessage: statusMessage,
                    isEmpty: polls.isEmpty,
                    emptyText: "No polls yet. Create or join one!"
                )
                ScrollView {
                    VStack(spacing: 12) {
                        ForEach(Array(polls.enumerated()), id: \.element.id) { index, poll in
                            NavigationLink(destination: PollDetailView(poll: poll)) {
                                PollTile(poll: poll, colorIndex: index)
                            }
                        }
                    }
                    .frame(maxWidth: .infinity, alignment: .center)
                    .padding(.vertical, 8)
                }
                PollActionButtons(
                    showAddPoll: $showAddPoll,
                    showJoinPoll: $showJoinPoll,
                    newPollName: $newPollName,
                    isCreatingPoll: $isCreatingPoll,
                    createPoll: createPoll,
                    inviteCode: $inviteCode,
                    isJoiningPoll: $isJoiningPoll,
                    joinPoll: joinPoll
                )
                .frame(maxWidth: .infinity, alignment: .center)
            }
            .padding()
            .navigationTitle("BiteBattle")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .navigationBarTrailing) {
                    AccountButton()
                }
            }
            .onAppear(perform: fetchPolls)
        }
        .sheet(isPresented: $showAddPoll) {
            AppBackground {
                VStack(spacing: 16) {
                    Text("New Poll Name")
                        .font(.headline)
                        .foregroundColor(AppColors.textPrimary)
                    AppTextField(placeholder: "Poll Name", text: $newPollName)
                    AppButton(
                        title: isCreatingPoll ? "Creating..." : "Create",
                        isLoading: isCreatingPoll,
                        isDisabled: newPollName.isEmpty,
                        action: createPoll
                    )
                    Button("Cancel") { showAddPoll = false }
                        .foregroundColor(.red)
                }
                .padding()
            }
            .interactiveDismissDisabled(isCreatingPoll)
        }
        .sheet(isPresented: $showJoinPoll) {
            AppBackground {
                VStack(spacing: 16) {
                    Text("Enter Invite Code")
                        .font(.headline)
                        .foregroundColor(AppColors.textPrimary)
                    AppTextField(placeholder: "Invite Code", text: $inviteCode)
                    AppButton(
                        title: isJoiningPoll ? "Joining..." : "Join",
                        isLoading: isJoiningPoll,
                        isDisabled: inviteCode.isEmpty,
                        action: joinPoll
                    )
                    Button("Cancel") { showJoinPoll = false }
                        .foregroundColor(.red)
                }
                .padding()
            }
            .interactiveDismissDisabled(isJoiningPoll)
        }
    }

    func fetchPolls() {
        isLoading = true
        APIClient.shared.fetchPolls { result in
            DispatchQueue.main.async {
                isLoading = false
                switch result {
                case .success(let polls):
                    self.polls = polls.sorted { $0.updated_at > $1.updated_at }
                    self.statusMessage = nil
                case .failure(let error):
                    self.statusMessage = error.localizedDescription
                }
            }
        }
    }

    func createPoll() {
        guard !newPollName.isEmpty else { return }
        isCreatingPoll = true
        APIClient.shared.createPoll(name: newPollName) { result in
            DispatchQueue.main.async {
                isCreatingPoll = false
                switch result {
                case .success(_):
                    newPollName = ""
                    // Dismiss the sheet, then fetch polls
                    DispatchQueue.main.asyncAfter(deadline: .now() + 0.1) {
                        showAddPoll = false
                        fetchPolls()
                    }
                case .failure(let error):
                    statusMessage = error.localizedDescription
                }
            }
        }
    }

    func joinPoll() {
        guard !inviteCode.isEmpty else { return }
        isJoiningPoll = true
        // Fetch pollId by invite code (simulate or add a lookup if needed)
        APIClient.shared.fetchPolls { result in
            DispatchQueue.main.async {
                switch result {
                case .success(let polls):
                    if let poll = polls.first(where: { $0.invite_code == inviteCode }) {
                        APIClient.shared.joinPoll(pollId: poll.id, inviteCode: inviteCode) { joinResult in
                            DispatchQueue.main.async {
                                isJoiningPoll = false
                                switch joinResult {
                                case .success(_):
                                    inviteCode = ""
                                    DispatchQueue.main.asyncAfter(deadline: .now() + 0.1) {
                                        showJoinPoll = false
                                        fetchPolls()
                                    }
                                case .failure(let error):
                                    statusMessage = error.localizedDescription
                                }
                            }
                        }
                    } else {
                        isJoiningPoll = false
                        statusMessage = "No poll found for invite code."
                    }
                case .failure(let error):
                    isJoiningPoll = false
                    statusMessage = error.localizedDescription
                }
            }
        }
    }

    struct GradientBackground<Content: View>: View {
        let content: Content
        init(@ViewBuilder content: () -> Content) {
            self.content = content()
        }
        var body: some View {
            ZStack {
                LinearGradient(
                    gradient: Gradient(colors: [Color.pink.opacity(0.7), Color.orange.opacity(0.7)]),
                    startPoint: .topLeading,
                    endPoint: .bottomTrailing
                )
                .ignoresSafeArea()
                content
            }
        }
    }

    struct AppIcon: View {
        var body: some View {
            Image(systemName: "fork.knife.circle.fill")
                .resizable()
                .scaledToFit()
                .frame(width: 60, height: 60)
                .foregroundColor(.pink)
                .shadow(radius: 6)
        }
    }

    struct TitleText: View {
        let title: String
        var body: some View {
            Text(title)
                .font(.title)
                .fontWeight(.bold)
                .foregroundColor(AppColors.textPrimary)
        }
    }

    struct StatusOrLoadingView: View {
        let isLoading: Bool
        let statusMessage: String?
        let isEmpty: Bool
        let emptyText: String
        var body: some View {
            if isLoading {
                ProgressView()
            } else if let status = statusMessage {
                Text(status)
                    .foregroundColor(.red)
            } else if isEmpty {
                Text(emptyText)
                    .foregroundColor(.gray)
            }
        }
    }

    struct PollTile: View {
        let poll: Poll
        let colorIndex: Int

        // Cream color palette
        private let backgroundColors: [Color] = [
            Color(hex: "#fff7e3"),
            Color(hex: "#fff7e4"),
            Color(hex: "#fff8e4"),
            Color(hex: "#fff8e3"),
            Color(hex: "#fff8e5"),
            Color(hex: "#ffffec"),
            Color(hex: "#fffae6"),
            Color(hex: "#fffce8"),
            Color(hex: "#fff9e5"),
        ]
        private let borderColor = Color(hex: "#ffa43d")

        var isOwner: Bool {
            poll.role == "owner"
        }

        // Format the created_at string to a date string, or show as-is if parsing fails
        var formattedDate: String {
            let input = poll.created_at // This is a String
            let isoFormatter = ISO8601DateFormatter()
            if let date = isoFormatter.date(from: input) {
                let formatter = DateFormatter()
                formatter.dateStyle = .medium
                return formatter.string(from: date)
            } else {
                return input // fallback: show the raw string
            }
        }

        var numberOfMembers: Int {
            poll.members.count
        }

        var body: some View {
            VStack(spacing: 8) {
                HStack {
                    Image(systemName: isOwner ? "crown.fill" : "person.2.fill")
                        .foregroundColor(isOwner ? borderColor : .gray)
                        .imageScale(.large)
                    Spacer()
                }
                .padding(.bottom, 2)

                Text(poll.name)
                    .font(.headline)
                    .fontWeight(.bold)
                    .foregroundColor(.black)
                    .multilineTextAlignment(.center)
                    .frame(maxWidth: .infinity)

                Text("By \(poll.created_by ?? "Unknown")")
                    .font(.subheadline)
                    .foregroundColor(.gray)
                    .multilineTextAlignment(.center)

                HStack {
                    Text(formattedDate)
                        .font(.caption)
                        .foregroundColor(.gray)
                    Spacer()
                    HStack(spacing: 4) {
                        Image(systemName: "person.3.fill")
                            .foregroundColor(borderColor)
                            .imageScale(.small)
                        Text("\(numberOfMembers)")
                            .font(.caption)
                            .foregroundColor(.gray)
                    }
                }
                .padding(.top, 4)
            }
            .padding()
            .background(
                RoundedRectangle(cornerRadius: 16)
                    .fill(backgroundColors[colorIndex % backgroundColors.count])
            )
            .overlay(
                RoundedRectangle(cornerRadius: 16)
                    .stroke(borderColor, lineWidth: 2)
            )
            .shadow(color: borderColor.opacity(0.08), radius: 6, x: 0, y: 3)
            .padding(.horizontal, 16)
            .frame(maxWidth: 400)
        }
    }

    struct PollActionButtons: View {
        @Binding var showAddPoll: Bool
        @Binding var showJoinPoll: Bool
        @Binding var newPollName: String
        @Binding var isCreatingPoll: Bool
        let createPoll: () -> Void
        @Binding var inviteCode: String
        @Binding var isJoiningPoll: Bool
        let joinPoll: () -> Void
        var body: some View {
            HStack(spacing: 16) {
                Button(action: { showAddPoll = true }) {
                    Text("Create Poll")
                        .fontWeight(.semibold)
                        .padding()
                        .background(Color.pink.opacity(0.9))
                        .foregroundColor(.white)
                        .cornerRadius(10)
                }
                Button(action: { showJoinPoll = true }) {
                    Text("Join Poll")
                        .fontWeight(.semibold)
                        .padding()
                        .background(Color.orange.opacity(0.9))
                        .foregroundColor(.white)
                        .cornerRadius(10)
                }
            }
            .padding(.top, 12)
        }
    }

    struct AccountButton: View {
        var body: some View {
            NavigationLink(destination: AccountView()) {
                HStack(spacing: 6) {
                    Image(systemName: "person.crop.circle")
                        .foregroundColor(.white)
                    Text("Account")
                        .foregroundColor(.white)
                }
                .padding(.vertical, 6)
                .padding(.horizontal, 10)
                .background(Color.pink.opacity(0.7))
                .cornerRadius(10)
            }
        }
    }
}


